package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/pkg/taskmaster"
)

var (
	confPathFlag   string
	socketPathFlag string
)

var grpcServer = grpc.NewServer()

func init() {
	log.SetPrefix("taskmasterd: ")

	flag.StringVar(&confPathFlag, "conf", "./taskmaster.yaml", "config file path")
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	conf, err := config.Parse(confPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	server, err := newTaskmasterServer(conf)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	pb.RegisterTaskmasterServer(grpcServer, server)

	l, err := net.Listen("unix", socketPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

type taskmasterServer struct {
	pb.UnimplementedTaskmasterServer

	conf   config.Conf
	master *taskmaster.Master
}

func newTaskmasterServer(conf config.Conf) (*taskmasterServer, error) {
	master, err := createMasterFromConfig(conf)
	if err != nil {
		return nil, err
	}

	master.AutoStart()

	return &taskmasterServer{
		conf:   conf,
		master: master,
	}, nil
}

var errTaskNotFound = errors.New("taskmaster: task not found")

func (s *taskmasterServer) ListTasks(ctx context.Context, _ *emptypb.Empty) (*pb.TasksList, error) {
	tasks := s.master.Tasks()
	result := make([]*pb.Task, len(tasks))

	for name, task := range tasks {
		result = append(result, taskToProto(name, task))
	}

	return &pb.TasksList{
		Tasks: result,
	}, nil
}

func (s *taskmasterServer) GetTask(ctx context.Context, taskIdentifier *pb.TaskIdentifier) (*pb.Task, error) {
	task, ok := s.master.Tasks()[taskIdentifier.Name]
	if !ok {
		return nil, errTaskNotFound
	}

	return taskToProto(taskIdentifier.Name, task), nil
}

func (s *taskmasterServer) StartTask(ctx context.Context, taskIdentifier *pb.TaskIdentifier) (*emptypb.Empty, error) {
	task, ok := s.master.Tasks()[taskIdentifier.Name]
	if !ok {
		return nil, errTaskNotFound
	}

	go task.Start()

	return &emptypb.Empty{}, nil
}

// func (s *taskmasterServer) RestartTask(ctx context.Context, taskIdentifier *pb.TaskIdentifier) (*emptypb.Empty, error) {
// 	return &emptypb.Empty{}, nil
// }

func (s *taskmasterServer) StopTask(ctx context.Context, taskIdentifier *pb.TaskIdentifier) (*emptypb.Empty, error) {
	task, ok := s.master.Tasks()[taskIdentifier.Name]
	if !ok {
		return nil, errTaskNotFound
	}

	go task.Stop()

	return &emptypb.Empty{}, nil
}

func (s *taskmasterServer) Reload(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	master, err := createMasterFromConfig(s.conf)
	if err != nil {
		return nil, err
	}

	go func() {
		s.master.Shutdown()
		s.master = master
		master.AutoStart()
	}()

	return &emptypb.Empty{}, nil
}

func (s *taskmasterServer) Stop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	go func() {
		grpcServer.GracefulStop()
		s.master.Shutdown()
	}()

	return &emptypb.Empty{}, nil
}

func createMasterFromConfig(conf config.Conf) (*taskmaster.Master, error) {
	master := taskmaster.NewMaster(log.New(os.Stdout, "taskmaster", 0))

	for name, program := range conf.Programs {
		taskAttr, err := programToTaskAttr(program)
		if err != nil {
			return nil, err
		}

		master.AddTask(name, program.NumProcs, taskAttr)
	}

	return master, nil
}

func programToTaskAttr(program config.Program) (taskmaster.TaskAttr, error) {
	argv := strings.Split(program.Cmd, " ")

	restartPolicy := taskmaster.RestartPolicyNum(program.AutoRestart)
	if restartPolicy == 0 {
		return taskmaster.TaskAttr{}, errors.New("unknown autorestart policy")
	}

	stopSignal := unix.SignalNum(program.StopSignal)
	if stopSignal == 0 {
		return taskmaster.TaskAttr{}, errors.New("unknown signal")
	}

	return taskmaster.TaskAttr{
		Bin:          argv[0],
		Argv:         argv,
		UMask:        program.UMask,
		Dir:          program.WorkingDir,
		AutoStart:    program.AutoStart,
		Restart:      restartPolicy,
		ExitCodes:    program.ExitCodes,
		StartRetries: program.StartRetries,
		StartTime:    time.Second * time.Duration(program.StartTime),
		StopSig:      stopSignal,
		StopTime:     time.Second * time.Duration(program.StopTime),
		Env:          envMapToSlice(program.Env),
	}, nil
}

func envMapToSlice(envMap map[string]string) []string {
	env := make([]string, len(envMap))
	i := 0

	for k, v := range envMap {
		env[i] = k + "=" + v
		i++
	}

	return env
}

func taskToProto(name string, task *taskmaster.Task) *pb.Task {
	t := &pb.Task{
		Name:   name,
		Status: statusToProto(task.Status()),
	}

	return t
}

func statusToProto(status taskmaster.Status) pb.Status {
	switch status {
	case taskmaster.StatusUnstarted:
		return pb.Status_UNSTARTED
	case taskmaster.StatusStarting:
		return pb.Status_STARTING
	case taskmaster.StatusRunning:
		return pb.Status_RUNNING
	case taskmaster.StatusStopping:
		return pb.Status_RUNNING
	case taskmaster.StatusStopped:
		return pb.Status_STOPPED
	case taskmaster.StatusErrored:
		return pb.Status_ERRORED
	default:
		return pb.Status_UNKNOWN
	}
}

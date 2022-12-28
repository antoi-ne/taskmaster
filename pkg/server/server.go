package server

import (
	"context"
	"log"
	"net"
	"syscall"

	"google.golang.org/grpc"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/pkg/taskmaster"
)

type server struct {
	pb.UnimplementedTaskmasterServer
}

func Run() error {
	m, err := taskmaster.NewMaster(log.Default())
	if err != nil {
		return err
	}

	m.AddTask("test", 1, taskmaster.TaskAttr{
		Bin:          "/bin/sleep",
		Argv:         []string{"sleep", "10"},
		AutoStart:    true,
		Restart:      taskmaster.RestartUnexpected,
		ExitCodes:    []int{0},
		StartRetries: 2,
		StartTime:    1,
		StopSig:      syscall.SIGABRT,
		StopTime:     3,
	})

	m.AutoStart()

	s := grpc.NewServer()

	pb.RegisterTaskmasterServer(s, &server{})

	l, err := net.Listen("unix", "/tmp/taskmaster.sock")
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *server) List(ctx context.Context, _ *pb.Empty) (*pb.ProgramDescList, error) {
	return &pb.ProgramDescList{}, nil
}

func (s *server) ProgramRestart(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) ProgramStart(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) ProgramStatus(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	return &pb.ProgramDesc{}, nil
}

func (s *server) ProgramStop(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) Reload(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) Stop(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

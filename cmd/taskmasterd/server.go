package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/pkg/taskmaster"
)

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

	return &taskmasterServer{
		conf:   conf,
		master: master,
	}, nil
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

package main

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
)

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

package main

import (
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/pkg/taskmaster"
)

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

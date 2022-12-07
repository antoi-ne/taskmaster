package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"pkg.coulon.dev/taskmaster/internal/manager"
	pb "pkg.coulon.dev/taskmaster/internal/proto"
)

type server struct {
	pb.UnimplementedTaskmasterServer
	m *manager.Manager
}

// Run opens a unix socket on the given socket and serves the internal/manager package through a gRPC service.
func Run(socket string, m *manager.Manager) error {
	s := grpc.NewServer()
	pb.RegisterTaskmasterServer(s, &server{
		m: m,
	})

	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *server) List(ctx context.Context, _ *pb.Empty) (*pb.ServiceStatusList, error) {
	var ss []*pb.ServiceStatus

	for n, s := range s.m.ListPrograms() {
		ss = append(ss, &pb.ServiceStatus{
			Name:   n,
			Status: pb.Status(s),
		})
	}

	return &pb.ServiceStatusList{
		Services: ss,
	}, nil
}

func (s *server) Reload(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	if err := s.m.Reload(); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) Restart(ctx context.Context, sv *pb.Service) (*pb.ServiceStatus, error) {
	if err := s.m.RestartProgram(sv.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(sv.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ServiceStatus{
		Name:   sv.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) Start(ctx context.Context, sv *pb.Service) (*pb.ServiceStatus, error) {
	if err := s.m.StartProgram(sv.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(sv.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ServiceStatus{
		Name:   sv.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) Status(ctx context.Context, sv *pb.Service) (*pb.ServiceStatus, error) {
	status, err := s.m.ProgramStatus(sv.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ServiceStatus{
		Name:   sv.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) Stop(ctx context.Context, sv *pb.Service) (*pb.ServiceStatus, error) {
	if err := s.m.StopProgram(sv.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(sv.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ServiceStatus{
		Name:   sv.Name,
		Status: pb.Status(status),
	}, nil
}

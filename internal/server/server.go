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
	s *grpc.Server
}

// Run opens a unix socket on the given socket and serves the internal/manager package through a gRPC service.
func Run(socket string, m *manager.Manager) error {
	s := grpc.NewServer()
	pb.RegisterTaskmasterServer(s, &server{
		m: m,
		s: s,
	})

	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *server) List(ctx context.Context, _ *pb.Empty) (*pb.ProgramDescList, error) {
	var ps []*pb.ProgramDesc

	for n, s := range s.m.ListPrograms() {
		ps = append(ps, &pb.ProgramDesc{
			Name:   n,
			Status: pb.Status(s),
		})
	}

	return &pb.ProgramDescList{
		Programs: ps,
	}, nil
}

func (s *server) Stop(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	s.m.StopAllAndWait()

	defer s.s.Stop()

	return &pb.Empty{}, nil
}

func (s *server) Reload(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	if err := s.m.Reload(); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ProgramRestart(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	if err := s.m.RestartProgram(p.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(p.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ProgramDesc{
		Name:   p.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) ProgramStart(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	if err := s.m.StartProgram(p.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(p.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ProgramDesc{
		Name:   p.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) ProgramStatus(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	status, err := s.m.ProgramStatus(p.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ProgramDesc{
		Name:   p.Name,
		Status: pb.Status(status),
	}, nil
}

func (s *server) ProgramStop(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	if err := s.m.StopProgram(p.Name); err != nil {
		return nil, err
	}

	status, err := s.m.ProgramStatus(p.Name)
	if err != nil {
		return nil, err
	}

	return &pb.ProgramDesc{
		Name:   p.Name,
		Status: pb.Status(status),
	}, nil
}

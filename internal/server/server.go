package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
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

	for n, p := range s.m.ListPrograms() {
		pd := &pb.ProgramDesc{
			Name:   n,
			Status: pb.Status(p.Status()),
		}

		pid, ok := p.Pid()
		if ok {
			pid := int32(pid)
			pd.Pid = &pid
		}

		ec, ok := p.ExitCode()
		if ok {
			ec := int32(ec)
			pd.Exitcode = &ec
		}

		ps = append(ps, pd)
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

func (s *server) ProgramRestart(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	if err := s.m.RestartProgram(p.Name); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ProgramStart(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	if err := s.m.StartProgram(p.Name); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ProgramStatus(ctx context.Context, p *pb.Program) (*pb.ProgramDesc, error) {
	pr, err := s.m.ProgramGet(p.Name)
	if err != nil {
		return nil, err
	}

	pd := &pb.ProgramDesc{
		Name:   p.Name,
		Status: pb.Status(pr.Status()),
	}

	pid, ok := pr.Pid()
	if ok {
		pid := int32(pid)
		pd.Pid = &pid
	}

	ec, ok := pr.ExitCode()
	if ok {
		ec := int32(ec)
		pd.Exitcode = &ec
	}

	ut, ok := pr.Uptime()
	if ok {
		pd.Uptime = durationpb.New(ut)
	}

	return pd, nil
}

func (s *server) ProgramStop(ctx context.Context, p *pb.Program) (*pb.Empty, error) {
	if err := s.m.StopProgram(p.Name); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

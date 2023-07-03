package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/internal/config"
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

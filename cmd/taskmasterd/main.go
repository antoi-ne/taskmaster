package main

import (
	"flag"
	"fmt"
	"log"

	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/internal/manager"
	"pkg.coulon.dev/taskmaster/internal/server"
)

var (
	confPathFlag   string
	socketPathFlag string
)

func init() {
	flag.StringVar(&confPathFlag, "conf", "./taskmaster.yaml", "config file path")
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	cf, err := config.Parse(confPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fmt.Printf("progs: %d\n", len(cf.Programs))

	m, err := manager.New(cf)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	m.AutoStart()

	if err := server.Run(socketPathFlag, m); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

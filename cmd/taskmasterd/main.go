package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pkg.coulon.dev/taskmaster/internal/manager"
	"pkg.coulon.dev/taskmaster/internal/server"
)

var (
	confPathFlag   string
	socketPathFlag string
)

func init() {
	log.SetPrefix("taskmasterd: ")

	flag.StringVar(&confPathFlag, "conf", "./taskmaster.yaml", "config file path")
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	m, err := manager.New(confPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	m.AutoStart()

	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		<-sigs

		log.Printf("SIGHUP received, reloading taskmasterd")

		m.Reload()
	}()

	if err := server.Run(socketPathFlag, m); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

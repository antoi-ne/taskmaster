package main

import (
	"flag"
	"log"
	"os/signal"
	"syscall"

	"pkg.coulon.dev/taskmaster/shell"
)

var (
	socketPathFlag string
)

func init() {
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	s := shell.New("tm>")

	signal.Notify(s.SignalChan(), syscall.SIGINT)

	defer s.Restore()
	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

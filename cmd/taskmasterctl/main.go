package main

import (
	"flag"
	"log"

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

	s := shell.New("tm>", shellHandler)

	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

func shellHandler(q *shell.Query) error {
	if len(q.Argv()) < 1 {
		return nil
	}

	switch q.Argv()[0] {
	case "exit":
		q.Exit()
	default:
		q.Println("unknown command")
	}

	return nil
}

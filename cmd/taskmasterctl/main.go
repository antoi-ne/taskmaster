package main

import (
	"context"
	"flag"
	"log"

	"pkg.coulon.dev/taskmaster/internal/client"
	"pkg.coulon.dev/taskmaster/internal/proto"
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

	c, err := client.Dial(socketPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer c.Close()

	s := shell.New("tm>", func(q *shell.Query) error {
		switch q.Argv()[0] {
		case "exit":
			q.Exit()
		case "list":
			sl, err := c.List(context.Background(), &proto.Empty{})
			if err != nil {
				return err
			}
			for _, s := range sl.Services {
				q.Println(s.Name + ": " + s.Status.String())
			}
		default:
			q.Println("unknown command")
		}
		return nil
	})

	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

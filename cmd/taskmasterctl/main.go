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
			pl, err := c.List(context.Background(), &proto.Empty{})
			if err != nil {
				return err
			}
			for _, p := range pl.Programs {
				q.Println(p.Name + ": " + p.Status.String())
			}
		case "reload":
			_, err := c.Reload(context.Background(), &proto.Empty{})
			if err != nil {
				return err
			}
			q.Println("taskmasterd reloaded.")
		case "stop":
			_, err := c.Stop(context.Background(), &proto.Empty{})
			if err != nil {
				return err
			}
			q.Println("server stopped. exiting.")
			q.Exit()
		default:
			q.Println("unknown command.")
		}
		return nil
	})

	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

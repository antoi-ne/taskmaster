package main

import (
	"context"
	"flag"
	"log"

	"pkg.coulon.dev/taskmaster/internal/client"
	pb "pkg.coulon.dev/taskmaster/internal/proto"
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
			pl, err := c.List(context.Background(), &pb.Empty{})
			if err != nil {
				return err
			}
			for _, p := range pl.Programs {
				q.Printf("%s: %s\n", p.Name, p.Status.String())
			}
		case "reload":
			_, err := c.Reload(context.Background(), &pb.Empty{})
			if err != nil {
				return err
			}
			q.Println("taskmasterd reloaded.")
		case "stop":
			_, err := c.Stop(context.Background(), &pb.Empty{})
			if err != nil {
				return err
			}
			q.Println("server stopped. exiting.")
			q.Exit()
		case "service":
			if len(q.Argv()) != 3 {
				q.Println("invalid syntax.")
				return nil
			}

			switch q.Argv()[2] {
			case "status":
				p, err := c.ProgramStatus(context.Background(), &pb.Program{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

				q.Printf("%s: %s\n", p.Name, p.Status.String())
				if p.Pid != nil {
					q.Printf("pid: %d\n", p.GetPid())
				}
				if p.Exitcode != nil {
					q.Printf("exit code: %d\n", p.GetExitcode())
				}

			case "start":
				_, err := c.ProgramStart(context.Background(), &pb.Program{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

			case "restart":
				_, err := c.ProgramRestart(context.Background(), &pb.Program{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

			case "stop":
				_, err := c.ProgramStop(context.Background(), &pb.Program{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

			default:
				q.Println("unknown subcommand.")
			}

		default:
			q.Println("unknown command.")
		}
		return nil
	})

	if err := s.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

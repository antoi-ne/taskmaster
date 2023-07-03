package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
	"pkg.coulon.dev/taskmaster/pkg/shell"
)

var (
	socketPathFlag string
)

func init() {
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	socketURL, err := url.JoinPath("unix://", socketPathFlag)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	conn, err := grpc.Dial(socketURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer conn.Close()

	client := pb.NewTaskmasterClient(conn)

	s := shell.New("tm>", func(q *shell.Query) error {
		switch q.Argv()[0] {
		case "help":
			q.Println("Available commands:")
			q.Println("    help                                      : list available commands")
			q.Println("    exit                                      : exit the shell")
			q.Println("    list                                      : list all services and their status")
			q.Println("    reload                                    : restart server with updated config file")
			q.Println("    stop                                      : stop the server and exit")
			q.Println("    service NAME [status/start/stop/restart]  : perform action on individual program")
		case "exit":
			q.Exit()
		case "reload":
			_, err := client.Reload(context.Background(), &emptypb.Empty{})
			if err != nil {
				return err
			}
			q.Println("taskmasterd reloaded.")
		case "stop":
			_, err := client.Stop(context.Background(), &emptypb.Empty{})
			if err != nil {
				return err
			}
			q.Println("server stopped. exiting.")
			q.Exit()
		case "list":
			tasksList, err := client.ListTasks(context.Background(), &emptypb.Empty{})
			if err != nil {
				return err
			}
			for _, p := range tasksList.Tasks {
				q.Printf("%s: %s\n", p.Name, p.Status.String())
			}
		case "service":
			if len(q.Argv()) != 3 {
				q.Println("invalid syntax.")
				return nil
			}

			switch q.Argv()[2] {
			case "status":
				p, err := client.GetTask(context.Background(), &pb.TaskIdentifier{
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
				if p.Uptime != nil {
					q.Printf("uptime: %s\n", p.GetUptime().AsDuration().Round(time.Second).String())
				}
				if p.Exitcode != nil {
					q.Printf("exit code: %d\n", p.GetExitcode())
				}

			case "start":
				_, err := client.StartTask(context.Background(), &pb.TaskIdentifier{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

			case "restart":
				_, err := client.RestartTask(context.Background(), &pb.TaskIdentifier{
					Name: q.Argv()[1],
				})
				if err != nil {
					q.Println(err.Error())
					return nil
				}

			case "stop":
				_, err := client.StopTask(context.Background(), &pb.TaskIdentifier{
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

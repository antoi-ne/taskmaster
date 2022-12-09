package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "pkg.coulon.dev/taskmaster/api/taskmasterpb"
)

type Client struct {
	pb.TaskmasterClient
	conn *grpc.ClientConn
}

func Dial(socket string) (*Client, error) {
	c := new(Client)

	conn, err := grpc.Dial("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conn = conn

	c.TaskmasterClient = pb.NewTaskmasterClient(conn)

	return c, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

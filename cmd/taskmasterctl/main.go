package main

import (
	"flag"

	"pkg.coulon.dev/taskmaster/shell"
)

var (
	confPathFlag   string
	socketPathFlag string
)

func init() {
	flag.StringVar(&socketPathFlag, "socket", "/tmp/taskmaster.sock", "server socket path")
}

func main() {
	flag.Parse()

	s := shell.New("tm>")

	s.HandleSignals()
}

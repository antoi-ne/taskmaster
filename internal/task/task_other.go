// go:build !linux

package task

import (
	"fmt"
	"os"
	"time"
)

func start(a Attr, exitChan chan int) (*Task, error) {
	t := new(Task)

	t.Attr = a
	t.exitChan = exitChan

	fds, err := t.createChildFds()
	if err != nil {
		return nil, fmt.Errorf("could not open fds on /dev/null: %w", err)
	}

	p, err := os.StartProcess(t.Bin, t.Argv, &os.ProcAttr{
		Dir:   t.Dir,
		Env:   t.Env,
		Files: fds,
	})
	if err != nil {
		return nil, fmt.Errorf("task process could not be started: %w", err)
	}
	t.proc = p

	t.startTime = time.Now()

	go t.monitor()

	return t, nil
}

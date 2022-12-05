package task

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// Attr holds attributes that will be apllied to a new Task
type Attr struct {
	// Path of the command's binary to execute
	Bin string
	// Command arguments
	Argv []string
	// Working directory of the Task.
	Dir string
	// Environment variables for the new Task.
	Env []string
	// Stdout and Stderr specify the Task's standard output and error files.
	Stdout *os.File
	Stderr *os.File
}

// Task stores information about a task and its related process.
type Task struct {
	Attr
	exitChan chan int
	proc     *os.Process
	mu       sync.RWMutex
	exited   bool
	exitCode int
	exitPid  int
}

// Start a new process and monitor its status.
func Start(a Attr, exitChan chan int) (*Task, error) {
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

	go t.monitor()

	return t, nil
}

// Pid returns the process ID of the task's underlying process.
func (t *Task) Pid() int {
	if t.Running() {
		return t.proc.Pid
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.exitPid
}

// Running returns true if the process has not yet exited, false otherwise.
func (t *Task) Running() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return !t.exited
}

// ExitCode returns the exit code of the process and panics if the process is still running.
func (t *Task) ExitCode() int {
	if t.Running() {
		panic("ExitCode called on Task which is still running")
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.exitCode
}

func (t *Task) Signal(sig os.Signal) error {
	if t.Running() {
		return t.proc.Signal(sig)
	} else {
		return errors.New("cannot send signal to proccess which is not running")
	}
}

// Kill sends a SIGKILL signal to the underlying process
func (t *Task) Kill() error {
	if t.Running() {
		return t.proc.Kill()
	} else {
		return errors.New("cannot kill proccess which is not running")
	}
}

// monitor waits for the process to exit, then saves its exit code and notify exitChan
func (t *Task) monitor() {
	ps, err := t.proc.Wait()
	if err != nil {
		panic(err)
	}

	t.mu.Lock()
	t.exited = true
	t.exitCode = ps.ExitCode()
	t.exitPid = ps.Pid()
	t.mu.Unlock()

	if t.exitChan != nil {
		t.exitChan <- t.exitCode
	}
}

func (t *Task) createChildFds() (fds []*os.File, err error) {
	fds = make([]*os.File, 3)

	// Open /dev/null in readonly mode for stdin.
	fds[0], err = os.Open(os.DevNull)
	if err != nil {
		return nil, err
	}

	// If attr.Stdout does not exist, open /dev/null in writeonly mode.
	fds[1] = t.Stdout
	if fds[1] == nil {
		fds[1], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	// If attr.Stderr does not exist, open /dev/null in writeonly mode.
	fds[2] = t.Stderr
	if fds[2] == nil {
		fds[2], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	return fds, nil
}

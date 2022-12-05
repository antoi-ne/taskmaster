package task

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

// TaskAttr holds attributes that will be apllied to a new Task
type TaskAttr struct {
	// Working directory of the Task.
	Dir string
	// Environment variables for the new Task.
	Env []string
	// Exit codes considered as a success.
	SuccessCodes []int
	// Signal which should be used to kill the process.
	KillSig os.Signal
	// Stdout and Stderr specify the Task's standard output and error files.
	Stdout *os.File
	Stderr *os.File
	// Channel which receives the task's exit code on exit
	ExitChan chan int
}

// Task stores information about a task and its related process.
type Task struct {
	proc *os.Process

	successCodes []int
	killSig      os.Signal
	exitChan     chan int

	mu       sync.RWMutex
	exited   bool
	exitCode int
	exitPid  int
}

// New starts a new process and monitors its status.
func New(name string, argv []string, attr *TaskAttr) (*Task, error) {
	t := new(Task)

	if attr.SuccessCodes != nil {
		t.successCodes = attr.SuccessCodes
	} else {
		t.successCodes = []int{0}
	}

	if t.killSig != nil {
		t.killSig = attr.KillSig
	} else {
		t.killSig = syscall.SIGKILL
	}

	t.exitChan = attr.ExitChan

	fds, err := attr.createChildFds()
	if err != nil {
		return nil, fmt.Errorf("could not open fds on /dev/null: %w", err)
	}

	p, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir:   attr.Dir,
		Env:   attr.Env,
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

// Success returns true if the process exited with one of the specified exit codes and panics if the process is still running.
func (t *Task) Success() bool {
	x := t.ExitCode()

	for _, c := range t.successCodes {
		if x == c {
			return true
		}
	}
	return false
}

// Kill sends the predefined kill signal to the task's process.
func (t *Task) Kill() error {
	return t.proc.Signal(t.killSig)
}

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

func (attr *TaskAttr) createChildFds() (fds []*os.File, err error) {
	fds = make([]*os.File, 3)

	// Open /dev/null in readonly mode for stdin.
	fds[0], err = os.Open(os.DevNull)
	if err != nil {
		return nil, err
	}

	// If attr.Stdout does not exist, open /dev/null in writeonly mode.
	fds[1] = attr.Stdout
	if fds[1] == nil {
		fds[1], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	// If attr.Stderr does not exist, open /dev/null in writeonly mode.
	fds[2] = attr.Stderr
	if fds[2] == nil {
		fds[2], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	return fds, nil
}

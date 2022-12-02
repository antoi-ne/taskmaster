package task

import (
	"fmt"
	"os"
	"sync"
)

// TaskAttr holds attributes that will be apllied to a new Task
type TaskAttr struct {
	// Working directory of the Task.
	Dir string
	// Environment variables for the new Task.
	Env map[string]string
	// Exit codes considered as a success.
	SuccessCodes []int
	// Signal which should be used to kill the process.
	KillSig os.Signal
	// Stdout and Stderr specify the Task's standard output and error files.
	Stdout *os.File
	Stderr *os.File
}

// Task stores information about a task and its related process.
type Task struct {
	proc *os.Process

	successCodes []int
	killSig      os.Signal

	mu       sync.RWMutex
	exited   bool
	exitCode int
}

// New starts a new process and monitors its status.
func New(name string, argv []string, attr *TaskAttr) (*Task, error) {
	t := new(Task)

	t.successCodes = attr.SuccessCodes

	fds, err := attr.createChildFds()
	if err != nil {
		return nil, err
	}

	p, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir:   attr.Dir,
		Env:   envMapToSlice(attr.Env),
		Files: fds,
	})
	if err != nil {
		return nil, fmt.Errorf("task process could not be started: %w", err)
	}
	t.proc = p

	go t.monitor()

	return t, nil
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

// Kill sends the predefined kill signal to the task's process
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
	t.mu.Unlock()
}

func (attr *TaskAttr) createChildFds() (fds []*os.File, err error) {
	fds = make([]*os.File, 3)

	// Open /dev/null in readonly mode for stdin
	fds[0], err = os.Open(os.DevNull)
	if err != nil {
		return nil, err
	}

	// If attr.Stdout does not exist, open /dev/null in writeonly mode
	fds[1] = attr.Stdout
	if fds[1] == nil {
		fds[1], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	// If attr.Stderr does not exist, open /dev/null in writeonly mode
	fds[2] = attr.Stderr
	if fds[2] == nil {
		fds[2], err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
	}

	return fds, nil
}

func envMapToSlice(envMap map[string]string) []string {
	env := make([]string, len(envMap))
	i := 0

	for k, v := range envMap {
		env[i] = k + "=" + v
		i++
	}

	return env
}

package process

import (
	"os"
	"sync"
	"time"
)

// Process represents a child process whose state is monitored.
type Process struct {
	bin  string
	argv []string

	proc      *os.Process
	exitCh    chan int
	startTime time.Time

	mx     sync.RWMutex
	ps     *os.ProcessState // ps = nil means the process is still running.
	uptime time.Duration    // not set until the process exits.
}

// Start creates a new child process and starts a monitoring routine to keep track of its state.
func Start(bin string, argv []string) (p *Process, err error) {
	p = &Process{
		bin:       bin,
		argv:      argv,
		exitCh:    make(chan int, 1),
		startTime: time.Now(),
	}

	p.proc, err = os.StartProcess(bin, argv, &os.ProcAttr{})
	if err != nil {
		return
	}

	// Start the monitoring routine.
	go p.monitor()

	return
}

// C retuns a channel which will be notified of the process exit code on its termination.
func (p *Process) C() chan int {
	return p.exitCh
}

// monitor waits for the process to exit then sets the process state and notifies the exit channel.
func (p *Process) monitor() {
	ps, _ := p.proc.Wait()

	p.mx.Lock()
	p.ps = ps
	p.uptime = time.Since(p.startTime)
	p.mx.Unlock()

	p.exitCh <- ps.ExitCode()
}

// Running returns true if the process has not exited yet.
func (p *Process) Running() bool {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.ps == nil
}

// ExitCode returns the exit code of the process if it exited and panics if the process is still running.
func (p *Process) ExitCode() int {
	p.mx.RLock()
	defer p.mx.RUnlock()

	if p.ps == nil {
		panic("process.ExitCode called on running process")
	}

	return p.ps.ExitCode()
}

// Pid returns the ID assigned to the process.
func (p *Process) Pid() int {
	p.mx.RLock()
	defer p.mx.RUnlock()

	if p.ps == nil {
		return p.proc.Pid
	}

	return p.ps.Pid()
}

// Uptime returns the time since the process started.
// If the process already exited, Uptime returns the time until termination.
func (p *Process) Uptime() time.Duration {
	p.mx.RLock()
	defer p.mx.RUnlock()

	if p.ps == nil {
		return time.Since(p.startTime)
	}

	return p.uptime
}

// Signal sends the given signal to the process.
func (p *Process) Signal(sig os.Signal) error {
	return p.proc.Signal(sig)
}

// Kill forces the process to exit imediately.
func (p *Process) Kill() error {
	return p.proc.Kill()
}

package program

import (
	"os"
	"sync"
	"time"

	"pkg.coulon.dev/taskmaster/task"
)

// Attr contains the configuration of a program.
type Attr struct {
	Argv         []string
	Bin          string
	UMask        uint32
	Dir          string
	AutoStart    bool
	Restart      RestartPolicy
	ExitCodes    []int
	StartRetries uint
	StartTime    time.Duration
	StopSig      os.Signal
	StopTime     time.Duration
	Stdout       *os.File
	Stderr       *os.File
	Env          []string
}

// Program manages executions of a single command.
type Program struct {
	Attr
	ta         task.Attr
	task       *task.Task
	exitChan   chan int      // Channel which receives the exit code of the task when it exits.
	actLock    sync.Mutex    // Lock to prevent multiple actions to be performed at the same time (start/stop/restart).
	actChan    chan struct{} // receives a notification when any action is started (used to block the monitor goroutine until the action is finished).
	statusLock sync.RWMutex  // Lock to prevent status to be written to by multiple goroutines at the same time.
	status     Status        // Current status of the program.
}

// New creates and returns a new Program with the given attributes.
func New(a Attr) (*Program, error) {
	p := new(Program)

	p.Attr = a

	p.exitChan = make(chan int)
	p.actChan = make(chan struct{})

	p.status = StatusUnstarted

	p.ta = task.Attr{
		Bin:    a.Bin,
		Argv:   a.Argv,
		Dir:    a.Dir,
		Env:    a.Env,
		Stdout: a.Stdout,
		Stderr: a.Stderr,
	}

	go p.monitor()

	return p, nil
}

// Status returns the current state of the program.
func (p *Program) Status() Status {
	p.statusLock.RLock()
	defer p.statusLock.RUnlock()

	return p.status
}

// Start starts the underlying tasks of the program. Waits for the operation to be finished.
func (p *Program) Start() error {
	p.tryStart()

	return nil
}

// Stop kills the tasks of the program by sending a signal. Waits for the operation to be finished.
func (p *Program) Stop() error {
	p.tryStop()

	return nil
}

// Restart stops then starts the program. Waits for the operation to be finished.
func (p *Program) Restart() error {
	p.tryStop()
	p.tryStart()

	return nil
}

func (p *Program) isExitCodeExpected(code int) bool {
	for _, c := range p.ExitCodes {
		if code == c {
			return true
		}
	}

	return false
}

func (p *Program) setStatus(s Status) {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()

	p.status = s
}

// monitor tries starting the program, then it monitors the process for multiple possible cases:
// - if the proccess exits by itself, restart it or not depending on the restart policy;
// - if an instruction is sent (STOP/RESTART), execute it.
func (p *Program) monitor() {
	for {
		select {
		// If the task exits, restart or not depending on the restart policy.
		case ec := <-p.exitChan:
			p.applyRestartPolicy(ec)
		// If any action is started, wait until it is finished.
		case <-p.actChan:
			p.actLock.Lock()
			p.actLock.Unlock()
		}
	}
}

// tryStart will try starting the task, set the appropriate status and return true if it succedded.
func (p *Program) tryStart() {
	p.actLock.Lock()
	defer p.actLock.Unlock()

	// Tell the monitor goroutine to block until the action lock is unlocked
	p.actChan <- struct{}{}

	p.setStatus(StatusStarting)

start_loop:
	for i := uint(0); i < p.StartRetries; i++ {
		t, err := task.Start(p.ta, p.exitChan)
		if err != nil {
			continue start_loop
		}
		p.task = t

		select {
		// If the program exits before the end of startTime, try again.
		case <-p.exitChan:
			continue start_loop

		// If the process did not exit, go to next step.
		case <-time.After(p.StartTime):
			p.setStatus(StatusRunning)
			return
		}
	}

	p.setStatus(StatusErrored)
}

// tryStop will try starting the task with the defined stop signal. If it has not stopped before the end of StopTime, a KILLSIG will be sent to force the task to exit.
func (p *Program) tryStop() {
	p.actLock.Lock()
	defer p.actLock.Unlock()

	// Tell the monitor goroutine to block until the action lock is unlocked
	p.actChan <- struct{}{}

	p.setStatus(StatusStopping)

	p.task.Signal(p.StopSig)

	select {
	case <-p.exitChan:
		break
	case <-time.After(p.StopTime):
		p.task.Kill()
	}
	p.setStatus(StatusStopped)
	<-p.exitChan
}

func (p *Program) applyRestartPolicy(exitCode int) {
	p.actLock.Lock()
	defer p.actLock.Unlock()

	if p.isExitCodeExpected(exitCode) {
		p.setStatus(StatusStopped)
	} else {
		p.setStatus(StatusErrored)
	}

	switch p.Attr.Restart {
	case RestartNever:
		return

	case RestartUnexpected:
		if p.isExitCodeExpected(exitCode) {
			return
		}
		fallthrough

	case RestartAlways:
		p.tryStart()

	default:
		panic("unknown restart policy for program")
	}
}

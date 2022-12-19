package taskmaster

import (
	"log"
	"os"
	"sync"
	"time"

	"pkg.coulon.dev/taskmaster/pkg/process"
)

type TaskAttr struct {
	Bin          string
	Argv         []string
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

type Task struct {
	TaskAttr
	name   string
	logger *log.Logger

	proc        *process.Process
	stopMonitor chan struct{} // send signal to stop the monitoring function
	statusLock  sync.RWMutex  // mutex for the reading/writing on status
	status      Status
	actionLock  sync.Mutex // prevents running multiple actions at the same time (start/stop)
}

func NewTask(name string, logger *log.Logger, attr TaskAttr) (*Task, error) {
	t := &Task{
		name:        name,
		logger:      logger,
		TaskAttr:    attr,
		stopMonitor: make(chan struct{}),
		status:      StatusUnstarted,
	}

	return t, nil
}

func (t *Task) Start() (err error) {
	t.actionLock.Lock()
	defer t.actionLock.Unlock()

	t.logger.Printf("trying to start task %s\n", t.name)

	t.setStatus(StatusStarting)

	for i := uint(0); i < t.StartRetries; i++ {
		t.proc, err = process.Start(t.Bin, t.Argv)
		if err != nil {
			continue
		}

		select {
		case <-t.proc.C(): // if the process exits before startTime, retry
			t.logger.Printf("task %s exited before startTime, retrying\n", t.name)

			continue

		case <-time.After(t.StartTime):
			t.logger.Printf("task %s started.\n", t.name)

			t.setStatus(StatusRunning)

			go t.monitor()

			return nil
		}
	}

	t.logger.Printf("task %s could not be started.\n", t.name)

	t.setStatus(StatusErrored)

	return nil
}

func (t *Task) Stop() error {
	t.actionLock.Lock()
	defer t.actionLock.Unlock()

	t.logger.Printf("trying to stop task %s.\n", t.name)

	t.stopMonitor <- struct{}{}

	t.setStatus(StatusStopping)

	t.proc.Signal(t.StopSig)

	select {
	case <-t.proc.C():

	case <-time.After(t.StopTime):
		t.logger.Printf("could not stop task %s cleanly, force killing the task.\n", t.name)

		t.proc.Kill()
		<-t.proc.C()
	}

	t.logger.Printf("task %s stopped.\n", t.name)

	t.setStatus(StatusStopped)

	return nil
}

func (t *Task) Status() Status {
	t.statusLock.RLock()
	defer t.statusLock.RUnlock()

	return t.status
}

// monitor waits until the process exits then applies the restart policy.
// It can be exited by sending a signal to the stopMonitor channel when a restart is planned.
func (t *Task) monitor() {
	select {
	case ec := <-t.proc.C():
		t.applyRestartPolicy(ec)

	case <-t.stopMonitor:
	}

}

func (t *Task) applyRestartPolicy(ec int) {
	t.actionLock.Lock()
	defer t.actionLock.Unlock()

	switch t.Restart {
	case RestartUnexpected:

		if t.isExitCodeExpected(ec) {
			t.setStatus(StatusStopped)

			break
		}

		fallthrough

	case RestartAlways:
		t.logger.Printf("task %s has exited (code: %d). restarting because of restart policy.\n", t.name, ec)

		go t.Start()

		return
	}

	t.logger.Printf("task %s has exited (code: %d). not restarting because of restart policy.\n", t.name, ec)
}

func (t *Task) isExitCodeExpected(ec int) bool {
	for _, c := range t.ExitCodes {
		if ec == c {
			return true
		}
	}

	return false
}

func (t *Task) setStatus(s Status) {
	t.statusLock.Lock()
	defer t.statusLock.Unlock()

	t.status = s
}

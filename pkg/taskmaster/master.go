package taskmaster

import (
	"fmt"
	"log"
	"sync"
)

type Master struct {
	logger *log.Logger
	tasks  map[string]*Task
}

// NewMaster allocates a new task supervisor.
func NewMaster(logger *log.Logger) *Master {
	m := &Master{
		logger: logger,
		tasks:  make(map[string]*Task),
	}

	return m
}

// AddTask creates count new instances of a task to the master.
func (m *Master) AddTask(name string, count uint, attr TaskAttr) error {
	for i := uint(0); i < count; i++ {
		n := name
		if count > 1 {
			n = fmt.Sprintf("%d-%s", i, name)
		}

		t, err := NewTask(n, m.logger, attr)
		if err != nil {
			return err
		}

		m.tasks[n] = t
	}

	return nil
}

// AutoStart starts all task which have the autotart directive.
func (m *Master) AutoStart() {
	for _, t := range m.tasks {
		if t.attr.AutoStart {
			go t.Start()
		}
	}
}

// Shutdown stops all tasks and waits for all of them to exit.
func (m *Master) Shutdown() {
	var wg sync.WaitGroup

	for _, t := range m.tasks {
		switch t.Status() {
		case StatusStarting, StatusRunning:
			wg.Add(1)
			task := t

			go func() {
				task.Stop()
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func (m *Master) Tasks() map[string]*Task {
	return m.tasks
}

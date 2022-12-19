package taskmaster

import (
	"fmt"
	"log"
)

type Master struct {
	logger *log.Logger
	procs  map[string]*Task
}

func NewMaster(logger *log.Logger) (*Master, error) {
	m := &Master{
		logger: logger,
		procs:  make(map[string]*Task),
	}

	return m, nil
}

// AddTask creates count new instances of a task to the master.
func (m *Master) AddTask(name string, count int, attr TaskAttr) error {
	for i := int(0); i < count; i++ {
		n := name
		if count > 1 {
			n = fmt.Sprintf("%d-%s", i, name)
		}

		t, err := NewTask(n, m.logger, attr)
		if err != nil {
			return err
		}

		m.procs[n] = t
	}

	return nil
}

func (m *Master) AutoStart() {
	for n, t := range m.procs {
		if t.AutoStart {
			m.logger.Printf("task %s auto starting.\n", n)

			go t.Start()
		}
	}
}

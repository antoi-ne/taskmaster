package manager

import "pkg.coulon.dev/taskmaster/internal/config"

// Manager represents a group of programs which are supervised
type Manager struct {
	cfg *config.File

	progs map[string]*program
}

// New creates a new manager
func New(cf *config.File) (*Manager, error) {
	m := new(Manager)
	m.cfg = cf

	m.progs = make(map[string]*program)

	for n, p := range m.cfg.Programs {
		prog, err := newProgram(p)
		if err != nil {
			return nil, err
		}
		m.progs[n] = prog
	}

	return m, nil
}

func (m *Manager) Start() error {
	for _, p := range m.progs {
		if err := p.start(); err != nil {
			return err
		}
	}

	return nil
}

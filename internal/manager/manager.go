package manager

// Manager represents a group of Programs which are supervised
type Manager struct {
	progs []*Program
}

// New creates a new manager
func New() (*Manager, error) {
	m := new(Manager)

	m.progs = make([]*Program, 0)

	return m, nil
}

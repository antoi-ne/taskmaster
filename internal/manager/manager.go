package manager

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/program"
)

var ErrProgramNotFound = errors.New("program not found")

// Manager type contains multiple programs.
type Manager struct {
	progs map[string]*program.Program
}

// New creates a new manager from the given configuration.
func New(c *config.File) (*Manager, error) {
	m := new(Manager)

	m.progs = make(map[string]*program.Program)

	for n, p := range c.Programs {
		for i := uint(0); i < p.NumProcs; i++ {
			prog, err := createProgram(p)
			if err != nil {
				return nil, err
			}

			name := fmt.Sprintf("%s-%d", n, i)

			fmt.Println("new prog added")
			m.progs[name] = prog
		}
	}

	return m, nil
}

// AutoStart will try starting every program which is configured to start on launch (autostart).
func (m *Manager) AutoStart() {
	for _, p := range m.progs {
		if p.AutoStart {
			go p.Start()
		}
	}
}

func (m *Manager) ListPrograms() map[string]program.Status {
	l := make(map[string]program.Status)

	for n, p := range m.progs {
		l[n] = p.Status()
	}

	return l
}

func (m *Manager) ProgramStatus(name string) (program.Status, error) {
	p, ok := m.progs[name]
	if !ok {
		return 0, ErrProgramNotFound
	}

	return p.Status(), nil
}

func (m *Manager) StartProgram(name string) error {
	p, ok := m.progs[name]
	if !ok {
		return ErrProgramNotFound
	}

	switch p.Status() {
	case program.StatusUnstarted, program.StatusStopped, program.StatusErrored:
		go p.Start()
	default:
		return errors.New("the program is already running")
	}

	return nil
}

func (m *Manager) StopProgram(name string) error {
	p, ok := m.progs[name]
	if !ok {
		return ErrProgramNotFound
	}

	switch p.Status() {
	case program.StatusStarting, program.StatusRunning:
		go p.Stop()
	default:
		return errors.New("the program is already stopped")
	}

	return nil
}

func (m *Manager) RestartProgram(name string) error {
	if err := m.StopProgram(name); err != nil {
		return err
	}

	return m.StartProgram(name)
}

func createProgram(prog *config.Program) (*program.Program, error) {
	argv := strings.Fields(prog.Cmd)

	rp := program.RestartPolicyNum(prog.AutoRestart)
	if rp == 0 {
		return nil, errors.New("unknown auto restart policy")
	}

	ss := unix.SignalNum(prog.StopSignal)
	if ss == 0 {
		return nil, errors.New("unknown stop signal")
	}

	var of, ef *os.File
	var err error

	if prog.Stdout != "" {
		of, err = os.OpenFile(prog.Stdout, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return nil, err
		}
	}

	if prog.Stderr != "" {
		ef, err = os.OpenFile(prog.Stderr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return nil, err
		}
	}

	return program.New(program.Attr{
		Argv:         argv,
		Bin:          argv[0],
		UMask:        prog.UMask,
		Dir:          prog.WorkingDir,
		AutoStart:    prog.AutoStart,
		Restart:      rp,
		ExitCodes:    prog.ExitCodes,
		StartRetries: prog.StartRetries,
		StartTime:    time.Second * time.Duration(prog.StartTime),
		StopSig:      ss,
		StopTime:     time.Second * time.Duration(prog.StopTime),
		Stdout:       of,
		Stderr:       ef,
		Env:          envMapToSlice(prog.Env),
	})
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

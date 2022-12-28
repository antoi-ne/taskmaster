package manager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/internal/program"
)

var ErrProgramNotFound = errors.New("program not found")

// Manager type contains multiple programs.
type Manager struct {
	configPath string
	config     *config.File
	progs      map[string]*program.Program
}

// New creates a new manager from the given configuration file path.
func New(configPath string) (*Manager, error) {
	m := new(Manager)

	m.configPath = configPath

	conf, err := config.Parse(m.configPath)
	if err != nil {
		return nil, err
	}
	m.config = conf

	progs, err := loadConfigIntoPrograms(conf)
	if err != nil {
		return nil, err
	}

	m.progs = progs

	return m, nil
}

func loadConfigIntoPrograms(c *config.File) (map[string]*program.Program, error) {
	progs := make(map[string]*program.Program)

	for n, p := range c.Programs {
		for i := uint(0); i < p.NumProcs; i++ {
			prog, err := createProgram(p)
			if err != nil {
				return nil, err
			}

			name := n
			if p.NumProcs > 1 {
				name = fmt.Sprintf("%s-%d", n, i)
			}

			progs[name] = prog
		}
	}

	return progs, nil
}

// AutoStart will try starting every program which is configured to start on launch (autostart). Nonblocking.
func (m *Manager) AutoStart() {
	for n, p := range m.progs {

		log.Printf("autostarting %s", n)

		if p.AutoStart {
			go p.Start()
		}
	}
}

// StopAllAndWait stops all running programs then waits for all of them to be exited.
func (m *Manager) StopAllAndWait() {
	var wg sync.WaitGroup

	for _, p := range m.progs {
		p := p

		switch p.Status() {
		case program.StatusStarting, program.StatusRunning:
			wg.Add(1)

			go func() {
				p.Stop()
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

// Reload stops all running program, parses the config file and starts all programs with the autostart directive.
func (m *Manager) Reload() error {
	log.Printf("reloading taskmasterd")

	conf, err := config.Parse(m.configPath)
	if err != nil {
		return err
	}

	// TODO: only reload changed/added programs
	// var unchanged []string

	// for n, c := range conf.Programs {
	// 	orig, ok := m.config.Programs[n]
	// 	if !ok {
	// 		continue
	// 	}

	// 	if reflect.DeepEqual(c, orig) {
	// 		unchanged = append(unchanged, n)
	// 	}
	// }

	// newProgs := make(map[string]*program.Program)

	m.StopAllAndWait()

	// delete all programs
	for k := range m.progs {
		delete(m.progs, k)
	}

	progs, err := loadConfigIntoPrograms(conf)
	if err != nil {
		return err
	}

	m.progs = progs

	m.AutoStart()

	return nil
}

func (m *Manager) ListPrograms() map[string]*program.Program {
	l := make(map[string]*program.Program)

	for n, p := range m.progs {
		l[n] = p
	}

	return l
}

func (m *Manager) ProgramGet(name string) (*program.Program, error) {
	p, ok := m.progs[name]
	if !ok {
		return nil, ErrProgramNotFound
	}

	return p, nil
}

func (m *Manager) StartProgram(name string) error {
	p, ok := m.progs[name]
	if !ok {
		return ErrProgramNotFound
	}

	switch p.Status() {
	case program.StatusUnstarted, program.StatusStopped, program.StatusErrored:
		log.Printf("starting  %s", name)

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
		log.Printf("stopping  %s", name)

		go p.Stop()
	default:
		return errors.New("the program is already stopped")
	}

	return nil
}

func (m *Manager) RestartProgram(name string) error {
	// TODO

	return nil
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

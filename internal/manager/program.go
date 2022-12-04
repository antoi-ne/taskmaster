package manager

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/sys/unix"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/internal/task"
)

type program struct {
	cfg *config.Program

	bin     string
	argv    []string
	killSig os.Signal
	stdout  *os.File
	stderr  *os.File

	tasks []*task.Task
}

func newProgram(cp *config.Program) (*program, error) {
	p := new(program)
	p.cfg = cp

	p.argv = strings.Fields(p.cfg.Cmd)
	p.bin = p.argv[0]

	sig := unix.SignalNum(p.cfg.StopSignal)
	if sig == 0 {
		return nil, errors.New("unknown signal: " + p.cfg.StopSignal)
	}
	p.killSig = sig

	of, err := os.OpenFile(p.cfg.Stdout, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	p.stdout = of

	ef, err := os.OpenFile(p.cfg.Stderr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	p.stderr = ef

	p.tasks = make([]*task.Task, p.cfg.NumProcs)

	return p, nil
}

func (p *program) start() error {
	for i := 0; i < int(p.cfg.NumProcs); i++ {
		t, err := task.New(p.argv[0], p.argv, &task.TaskAttr{
			Dir:          p.cfg.WorkingDir,
			Env:          p.cfg.Env,
			SuccessCodes: p.cfg.ExitCode,
			KillSig:      p.killSig,
			Stdout:       p.stdout,
			Stderr:       p.stderr,
		})
		if err != nil {
			return err
		}
		p.tasks[i] = t
	}

	return nil
}

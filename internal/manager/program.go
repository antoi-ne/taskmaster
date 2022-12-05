package manager

import (
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	"pkg.coulon.dev/taskmaster/internal/config"
	"pkg.coulon.dev/taskmaster/internal/task"
)

type program struct {
	argv         []string
	bin          string
	count        uint
	umask        uint32
	dir          string
	autoStart    bool
	restart      RestartPolicy
	exitCodes    []int
	startRetries uint
	startTime    time.Duration
	stopSig      os.Signal
	stopTime     time.Duration
	stdout       *os.File
	stderr       *os.File
	env          []string

	tasks []*task.Task
}

func newProgram(cp *config.Program) (*program, error) {
	p := new(program)

	p.argv = strings.Fields(cp.Cmd)
	p.bin = p.argv[0]
	p.count = cp.NumProcs
	p.umask = cp.UMask
	p.dir = cp.WorkingDir
	p.autoStart = cp.AutoStart
	p.restart = RestartPolicy(cp.AutoRestart)
	if !p.restart.Exists() {
		return nil, errors.New("unknown autorestart policy")
	}
	p.exitCodes = cp.ExitCode
	p.startRetries = cp.StartRetries
	p.startTime = time.Second * time.Duration(cp.StartTime)
	sig := unix.SignalNum(cp.StopSignal)
	if sig == 0 {
		return nil, errors.New("unknown signal '" + cp.StopSignal + "'")
	}
	p.stopSig = sig
	p.stopTime = time.Second * time.Duration(cp.StopTime)
	of, err := os.OpenFile(cp.Stdout, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	p.stdout = of
	ef, err := os.OpenFile(cp.Stderr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	p.stderr = ef
	p.env = envMapToSlice(cp.Env)

	p.tasks = make([]*task.Task, p.count)

	return p, nil
}

func (p *program) start() error {
	for i := range p.tasks {
		t, err := task.New(p.bin, p.argv, &task.TaskAttr{
			Dir:          p.dir,
			Env:          p.env,
			SuccessCodes: p.exitCodes,
			KillSig:      p.stopSig,
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

func envMapToSlice(envMap map[string]string) []string {
	env := make([]string, len(envMap))
	i := 0

	for k, v := range envMap {
		env[i] = k + "=" + v
		i++
	}

	return env
}

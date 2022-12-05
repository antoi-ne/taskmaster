package program_test

import (
	"syscall"
	"testing"
	"time"

	"pkg.coulon.dev/taskmaster/internal/program"
)

func TestNew(t *testing.T) {
	p, err := program.New(program.Attr{
		Argv:         []string{"sleep", "3"},
		Bin:          "/bin/sleep",
		AutoStart:    true,
		Restart:      program.RestartUnexpected,
		ExitCodes:    []int{0},
		StartRetries: 3,
		StartTime:    1,
		StopSig:      syscall.SIGTERM,
		StopTime:     3,
		Env:          []string{"FOO=bar"},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("status initial: %s\n", p.Status().String())

	p.Start()
	time.Sleep(100000)

	t.Logf("status after start: %s\n", p.Status().String())

	time.Sleep(time.Second * 2)

	t.Logf("status after wait: %s\n", p.Status().String())

	time.Sleep(time.Second * 2)

	t.Logf("status after wait 2: %s\n", p.Status().String())
}

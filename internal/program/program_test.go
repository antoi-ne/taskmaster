package program_test

import (
	"syscall"
	"testing"

	"pkg.coulon.dev/taskmaster/internal/program"
)

func TestNew(t *testing.T) {
	_, err := program.New(program.Attr{
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
}

package process_test

import (
	"testing"

	"pkg.coulon.dev/taskmaster/pkg/process"
)

func TestProcess(t *testing.T) {
	p, err := process.Start("/bin/sleep", []string{"sleep", "1"})
	if err != nil {
		t.Fatal(err)
	}

	if p.Running() == false {
		t.Fatal("process exited too soon.")
	}

	if ret := <-p.C(); ret != 0 {
		t.Fatal("process did not return 0.")
	}

	if p.Running() == true {
		t.Fatal("process.Running returns true after exit.")
	}
}

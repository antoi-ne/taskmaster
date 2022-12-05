package task_test

import (
	"testing"
	"time"

	"pkg.coulon.dev/taskmaster/task"
)

func TestExitChan(t *testing.T) {

	ch := make(chan int, 1)

	tk, err := task.Start(task.Attr{
		Bin:  "/bin/test",
		Argv: []string{"test", "0"},
	}, ch)
	if err != nil {
		t.Fatal(err)
	}

	ret := <-ch

	if tk.Running() {
		t.Fatal("task is still running after ExitChan was notified")
	}

	if tk.ExitCode() != 0 {
		t.Fatal("task did not exit successfully")
	}

	if ret != tk.ExitCode() {
		t.Fatal("ExitChan code does not match with return value of ExitCode")
	}
}

func TestRunning(t *testing.T) {

	tk, err := task.Start(task.Attr{
		Bin:  "/bin/sleep",
		Argv: []string{"sleep", "1"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !tk.Running() {
		t.Fatal("task not running for long enough")
	}
	time.Sleep(time.Second * 2)
	if tk.Running() {
		t.Fatal("task should already be finished")
	}
	if tk.ExitCode() != 0 {
		t.Fatal("task did not return the expected exit code")
	}
}

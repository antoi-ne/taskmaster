package task_test

import (
	"testing"
	"time"

	"pkg.coulon.dev/taskmaster/internal/task"
)

func TestExitChan(t *testing.T) {

	ch := make(chan int, 1)

	tk, err := task.New("/bin/test", []string{"test", "0"}, &task.Attr{
		ExitChan: ch,
	})
	if err != nil {
		t.Fatal(err)
	}

	ret := <-ch

	if tk.Running() {
		t.Fatal("task is still running after ExitChan was notified")
	}

	if !tk.Success() {
		t.Fatal("task did not exit successfully")
	}

	if ret != tk.ExitCode() {
		t.Fatal("ExitChan code does not match with return value of ExitCode")
	}
}

func TestRunning(t *testing.T) {

	tk, err := task.New("/bin/sleep", []string{"sleep", "2"}, &task.Attr{
		SuccessCodes: []int{0},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !tk.Running() {
		t.Fatal("task not running for long enough")
	}
	time.Sleep(time.Second * 3)
	if tk.Running() {
		t.Fatal("task should already be finished")
	}
	if tk.ExitCode() != 0 {
		t.Fatal("task did not return the expected exit code")
	}
}

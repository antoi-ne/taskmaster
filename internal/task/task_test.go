package task_test

import (
	"testing"
	"time"

	"pkg.coulon.dev/taskmaster/internal/task"
)

func TestTask(t *testing.T) {

	tk, err := task.New("/bin/sleep", []string{"sleep", "2"}, &task.TaskAttr{
		SuccessCodes: []int{0},
		Env: map[string]string{
			"FOO": "bar",
		},
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

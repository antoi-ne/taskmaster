package program_test

import (
	"testing"

	"pkg.coulon.dev/taskmaster/internal/program"
)

func TestRestartPolicyExists(t *testing.T) {
	s1 := "never"
	s2 := "unexpected"
	s3 := "always"
	s4 := ""
	s5 := "fakepolicy"

	if !program.RestartPolicy(s1).Exists() {
		t.Fatal("Exists returned false for valid restart policy")
	}

	if !program.RestartPolicy(s2).Exists() {
		t.Fatal("Exists returned false for valid restart policy")
	}

	if !program.RestartPolicy(s3).Exists() {
		t.Fatal("Exists returned false for valid restart policy")
	}

	if program.RestartPolicy(s4).Exists() {
		t.Fatal("Exists returned true for an empty string")
	}

	if program.RestartPolicy(s5).Exists() {
		t.Fatal("Exists returned false for an invalid restart policy")
	}
}

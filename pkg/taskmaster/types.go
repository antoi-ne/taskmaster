package taskmaster

// RestartPolicy represents the autorestart configuration of a task.
type RestartPolicy int

const (
	RestartNever RestartPolicy = iota + 1
	RestartUnexpected
	RestartAlways
)

// String returs the string representation of a restart policy.
func (rp RestartPolicy) String() string {
	switch rp {
	case RestartNever:
		return "never"
	case RestartUnexpected:
		return "unexpected"
	case RestartAlways:
		return "always"
	default:
		return "unknown"
	}
}

// RestartPolicyNum returns the constant of the restart policy if it exists, 0 otherwise.
func RestartPolicyNum(name string) RestartPolicy {
	switch name {
	case "never":
		return RestartNever
	case "unexpected":
		return RestartUnexpected
	case "always":
		return RestartAlways
	default:
		return 0
	}
}

// Status represents all possible states of a program.
type Status int

const (
	StatusUnstarted Status = iota + 1
	StatusStarting
	StatusRunning
	StatusStopping
	StatusStopped
	StatusErrored
)

// String returns the string representation of a status.
func (s Status) String() string {
	switch s {
	case StatusUnstarted:
		return "UNSTARTED"
	case StatusStarting:
		return "STARTING"
	case StatusRunning:
		return "RUNNING"
	case StatusStopping:
		return "STOPPING"
	case StatusStopped:
		return "STOPPED"
	case StatusErrored:
		return "ERRORED"
	default:
		return "UNKNOWN"
	}
}

package program

// RestartPolicy represents the autorestart configuration of a program.
type RestartPolicy string

const (
	// RestartNever means the program will not restart automatically after a shutdown.
	RestartNever RestartPolicy = "never"
	// RestartUnexpected means the program will restart on shutdown only if the return code is not expected.
	RestartUnexpected RestartPolicy = "unexpected"
	// RestartAlways means the program will restart after every shutdown.
	RestartAlways RestartPolicy = "always"
)

func (rp RestartPolicy) Exists() bool {
	switch rp {
	case RestartNever:
		fallthrough
	case RestartUnexpected:
		fallthrough
	case RestartAlways:
		return true
	default:
		return false
	}
}

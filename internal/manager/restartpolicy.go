package manager

// RestartPolicy represents the autorestart configuration of a program
type RestartPolicy string

const (
	RestartNever      RestartPolicy = "never"      // RestartNever means the program will not restart automatically after a shutdown
	RestartUnexpected RestartPolicy = "unexpected" // RestartUnexpected means the program will restart on shutdown only if the return code is not expected
	RestartAlways     RestartPolicy = "always"     // RestartAlways means the program will restart after every shutdown

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

package envdetect

import (
	"os"
)

// RunningLocalDev returns true if the AMB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	s := os.Getenv("AMB_LOCAL")
	if len(s) > 0 {
		return true
	}

	return false
}

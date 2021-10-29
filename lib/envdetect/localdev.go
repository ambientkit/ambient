// Package envdetect provides flag based on environment variable being set.
package envdetect

import (
	"os"
	"strconv"
)

// RunningLocalDev returns true if the AMB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	result, _ := strconv.ParseBool(os.Getenv("AMB_LOCAL"))
	return result
}

// Package timezone allows setting the timezone for the application.
package timezone

import "os"

// Set the time zone based on the AMB_TIMEZONE environment variable or use
// EST time by default.
func Set() {
	// Get the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}

	os.Setenv("TZ", tz)
}

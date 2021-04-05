package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the time zone.
	SetTimezone()
}

func main() {
	// Generate a new private key for AES-256.
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	// Encode key in bytes to string for saving.
	key := hex.EncodeToString(bytes)
	fmt.Printf("AMB_SESSION_KEY=%v\n", key)
}

// SetTimezone the time zone based on the AMB_TIMEZONE environment variable or use
// EST time by default.
func SetTimezone() {
	// Get the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}

	os.Setenv("TZ", tz)
}

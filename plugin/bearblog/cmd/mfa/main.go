package main

import (
	"fmt"
	"log"
	"os"

	"github.com/josephspurrier/ambient/plugin/bearblog/lib/totp"
	"github.com/mdp/qrterminal/v3"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the time zone.
	SetTimezone()
}

func main() {
	username := os.Getenv("AMB_USERNAME")
	if len(username) == 0 {
		log.Fatalln("Environment variable missing:", "AMB_USERNAME")
	}

	issuer := os.Getenv("AMB_ISSUER")
	if len(issuer) == 0 {
		log.Fatalln("Environment variable missing:", "AMB_ISSUER")
	}

	// Generate a MFA.
	URI, secret, err := totp.GenerateURL(username, issuer)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Output the TOTP URI and config information.
	fmt.Printf("AMB_MFA_KEY=%v\n", secret)
	fmt.Println("")
	fmt.Println("Send this to a mobile phone to add it to an app like Google Authenticator or scan the QR code below:")
	fmt.Printf("%v\n", URI)

	config := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(URI, config)
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

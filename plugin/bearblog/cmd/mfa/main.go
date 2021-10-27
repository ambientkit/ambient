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
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) > 0 {
		os.Setenv("TZ", tz)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("Parameters missing:", "username issuer")
	}
	if len(os.Args) < 2 {
		log.Fatalln("Parameter missing:", "username")
	}
	if len(os.Args) < 3 {
		log.Fatalln("Parameter missing:", "issuer")
	}
	username := os.Args[1]
	issuer := os.Args[2]

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

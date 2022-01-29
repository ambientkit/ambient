package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/ambientkit/ambient/plugin/generic/bearblog/lib/passhash"
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
	if len(os.Args) != 2 {
		log.Fatalln("Incorrect number of arguments, you must pass in the password.")
	}

	// Generate a new private key.
	s, err := passhash.HashString(os.Args[1])
	if err != nil {
		log.Fatalln(err.Error())
	}

	sss := base64.StdEncoding.EncodeToString([]byte(s))
	fmt.Printf("AMB_PASSWORD_HASH=%v\n", sss)
}

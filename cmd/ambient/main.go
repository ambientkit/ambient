package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
)

func init() {
	// Set the time zone.
	tz := os.Getenv("AMB_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}

	os.Setenv("TZ", tz)
}

func main() {
	// Set up the application services.
	logger, mux, err := ambient.Boot(app.Plugins)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("web server listening on port: %v", port)
	logger.Fatal("", http.ListenAndServe(":"+port, mux))
}

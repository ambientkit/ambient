package main

import (
	"net/http"
	"os"

	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/app/lib/logger"
	"github.com/josephspurrier/ambient/app/lib/timezone"
	"github.com/sirupsen/logrus"
)

func init() {
	// Set the time zone.
	timezone.Set()
}

func main() {
	// Create the logger.
	l := logger.NewLogger("ambient", "1.0")
	//l.SetLevel(uint32(logrus.DebugLevel))
	l.SetLevel(uint32(logrus.InfoLevel))
	// l.SetLevel(logrus.ErrorLevel)
	// l.SetLevel(logrus.FatalLevel)

	// Set up the application services.
	mux, err := app.Boot(l)
	if err != nil {
		l.Fatal("", err.Error())
	}

	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	l.Info("web server listening on port: %v", port)
	l.Fatal("", http.ListenAndServe(":"+port, mux))
}

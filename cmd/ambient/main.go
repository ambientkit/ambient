package main

import (
	"log"
	"net/http"
	"os"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/app"
	"github.com/josephspurrier/ambient/lib/timezone"
)

func init() {
	// Set the time zone.
	timezone.Set()
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

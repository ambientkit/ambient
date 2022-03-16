package main

import (
	"log"
	"net/http"

	"github.com/ambientkit/ambient/internal/testutil"
)

func main() {
	_, pluginClient, h, err := testutil.Setup()
	if pluginClient != nil {
		defer pluginClient.Kill()
	}
	if err != nil {
		log.Fatalln(err.Error())
	}

	go http.ListenAndServe(":8080", h)

	select {}
}

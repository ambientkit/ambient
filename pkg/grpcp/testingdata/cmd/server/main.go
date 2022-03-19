package main

import (
	"log"

	"github.com/ambientkit/ambient/pkg/grpcp/testutil"
)

func main() {
	app, err := testutil.Setup(true)
	if err != nil {
		log.Fatalln(err.Error())
	}

	h, err := app.Handler()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Fatalln(app.ListenAndServe(h))
}

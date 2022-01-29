package internal

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/lib/requestclient"
)

var (
	// Globals.
	log ambient.AppLogger
	rc  *requestclient.RequestClient
)

// SetGlobals will set the variables used by the package.
func SetGlobals(l ambient.AppLogger, r *requestclient.RequestClient) {
	log = l
	rc = r
}

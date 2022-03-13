package grpcp

import (
	"fmt"
	"net/http"
)

// AmbientContextKey is a context key.
type AmbientContextKey string

const ambientRequestID AmbientContextKey = "ambient_requestid"

// requestID returns the request ID from the request context.
func requestID(r *http.Request) string {
	val := r.Context().Value(ambientRequestID)
	if val == nil {
		val = ""
	}
	return val.(string)
}

func pathkey(method string, path string) string {
	return fmt.Sprintf("%v %v", method, path)
}

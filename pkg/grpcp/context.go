package grpcp

import "fmt"

// AmbientContextKey is a context key.
type AmbientContextKey string

const ambientRequestID AmbientContextKey = "ambient_requestid"

func pathkey(method string, path string) string {
	return fmt.Sprintf("%v %v", method, path)
}

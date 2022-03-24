package hello

import (
	"context"
	"net/http"
)

// ContextKey is a context key.
type ContextKey string

const helloContextKey ContextKey = "hello_key"

// Set the value in the request context.
func Set(r *http.Request, value string) *http.Request {
	newContext := context.WithValue(r.Context(), helloContextKey, value)
	return r.WithContext(newContext)
}

// Get returns the request value from the request context.
func Get(r *http.Request) string {
	if r == nil {
		return ""
	}
	val := r.Context().Value(helloContextKey)
	if val == nil {
		val = ""
	}
	return val.(string)
}

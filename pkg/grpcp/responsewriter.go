package grpcp

import (
	"bytes"
	"net/http"
)

// ResponseWriter .
type ResponseWriter struct {
	statusCode int
	bytes      bytes.Buffer
	header     http.Header
}

// NewResponseWriter .
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		statusCode: 200,
		header:     make(http.Header),
	}
}

// Output .
func (w *ResponseWriter) Output() string {
	return w.bytes.String()
}

// StatusCode .
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

// Header .
func (w *ResponseWriter) Header() http.Header {
	return w.header
}

// Write .
func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.bytes.Write(b)
}

// WriteHeader .
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

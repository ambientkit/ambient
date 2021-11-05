package etagcache

import "net/http"

// CustomResponseWriter can be used
type CustomResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

// NewCustomResponseWriter stores the response without writing it.
func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		header: http.Header{},
	}
}

// Header return custom header.
func (w *CustomResponseWriter) Header() http.Header {
	return w.header
}

// Write to the custom body.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return len(w.body), nil
}

// WriteHeader will write to the custom status code.
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// WriteTo will write to the response writer.
func (w *CustomResponseWriter) WriteTo(wr http.ResponseWriter) {
	wr.WriteHeader(w.statusCode)
	for k, v := range w.header {
		for _, val := range v {
			wr.Header().Add(k, val)
		}
	}
	wr.Write(w.body)
}

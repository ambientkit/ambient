package requestuuid

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

// AmbientContextKey is a context key.
type AmbientContextKey string

// AmbientUUID is the context key to support unique requests.
const AmbientUUID AmbientContextKey = "ambient_requestid"

// Generate the request ID and store in the request context.
func Generate(r *http.Request) *http.Request {
	// Generate a unique request object, store the request for use by
	// Param(), then delete the request once the request is done to clean up.

	uuid := ""

	// if span := trace.SpanFromContext(r.Context()); span != nil {
	// 	//logger := l.clone()
	// 	//logger.ctx = ctx
	// 	// TODO: Determine if these need to be saved.
	// 	//span.SpanContext().TraceID()
	// 	//span.SpanContext().SpanID()
	// 	//l.Error("context found: %v", span.SpanContext().SpanID())
	// 	//return logger
	// 	if span.SpanContext().HasTraceID() {
	// 		uuid = span.SpanContext().TraceID().String()
	// 	}
	// }

	if len(uuid) == 0 {
		uuid, _ = UUID()
	}

	return Set(r, uuid)
}

// Set the request ID in the request context.
func Set(r *http.Request, value string) *http.Request {
	newContext := context.WithValue(r.Context(), AmbientUUID, value)

	// f := tracecontext.HTTPFormat{}
	// ot.
	// sc, ok := f.SpanContextFromRequest(r)
	// if ok {
	// 	newContext = trace.ContextWithRemoteSpanContext(newContext, sc)
	// }

	return r.WithContext(newContext)
}

// Get returns the request ID from the request context.
func Get(r *http.Request) string {
	if r == nil {
		return ""
	}
	val := r.Context().Value(AmbientUUID)
	if val == nil {
		val = ""
	}
	return val.(string)
}

// Middleware add a unique request ID to every request.
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := Generate(r)
		h.ServeHTTP(w, req)
	})
}

// UUID generates a UUID for use as an ID.
func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

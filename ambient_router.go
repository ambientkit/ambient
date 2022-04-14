package ambient

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AppRouter represents a router.
type AppRouter interface {
	Router

	ServeHTTP(w http.ResponseWriter, r *http.Request)
	SetNotFound(notFound http.Handler)
	SetServeHTTP(h func(w http.ResponseWriter, r *http.Request, err error))
}

// Router represents a router.
type Router interface {
	Handle(method string, path string, fn func(http.ResponseWriter, *http.Request) error)
	Get(path string, fn func(http.ResponseWriter, *http.Request) error)
	Post(path string, fn func(http.ResponseWriter, *http.Request) error)
	Patch(path string, fn func(http.ResponseWriter, *http.Request) error)
	Put(path string, fn func(http.ResponseWriter, *http.Request) error)
	Delete(path string, fn func(http.ResponseWriter, *http.Request) error)
	Head(path string, fn func(http.ResponseWriter, *http.Request) error)
	Options(path string, fn func(http.ResponseWriter, *http.Request) error)
	StatusError(status int, err error) error
	Error(status int, w http.ResponseWriter, r *http.Request)
	Param(r *http.Request, name string) string
	Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (err error)
}

// Route is a route for a router.
type Route struct {
	Method string
	Path   string
}

// CustomServeHTTP allows customization of error handling by the router.
type CustomServeHTTP func(log Logger, renderer Renderer,
	w http.ResponseWriter, r *http.Request, err error)

// SetupRouter sets the router with the NotFound handler and the default handler.
func SetupRouter(logger Logger, mux AppRouter, te Renderer, customServeHTTP CustomServeHTTP) {
	// Set the default handler.
	defaultServeHTTP := func(w http.ResponseWriter, r *http.Request, err error) {
		ctx, span := logger.Trace(r.Context(), "router: error handler")
		defer span.End()
		if err != nil {
			// Set default errors to internal server error.
			status := http.StatusInternalServerError
			friendlyError := "Darn, something went wrong."

			// If the error is a status error, then use the information.
			se, ok := err.(Error)
			if ok {
				if se.Status() > 0 {
					status = se.Status()
				}
				if len(se.Message()) > 0 {
					friendlyError = se.Message()
				}
			}

			span.SetAttributes(attribute.Int("http.status.code", status))
			span.SetAttributes(attribute.String("http.status.message", http.StatusText(status)))

			// Handle only errors.
			if status >= 400 {
				span.SetStatus(codes.Error, friendlyError)

				switch status {
				case 403:
					// Already logged on plugin access denials.
					friendlyError = "A plugin has been denied permission."
				case 404:
					// No need to log.
					friendlyError = "Darn, we cannot find the page."
				case 400:
					if err != nil {
						logger.Info("router error (%v): %v", status, err.Error())
					}
				case 500:
					if err != nil {
						logger.Error("router error (%v): %v", status, err.Error())
					}
				default:
					if err != nil {
						logger.Info("router error (%v): %v", status, err.Error())
					}
				}
				span.SetAttributes(attribute.String("http.err.friendly", http.StatusText(status)))

				if te != nil {
					err = te.Error(w, r, fmt.Sprintf("<h1>%v</h1>%v", status, friendlyError), status, nil, nil)
					if err != nil {
						if err != nil {
							logger.For(ctx).Info("router error in rendering error template (%v): %v", status, err.Error())
						}
						http.Error(w, "500 internal server error", http.StatusInternalServerError)
						return
					}
				} else {
					http.Error(w, friendlyError, status)
				}
			}
		}
	}

	// Use the custom handler if it's set.
	serveHTTP := defaultServeHTTP
	if customServeHTTP != nil {
		serveHTTP = func(w http.ResponseWriter, r *http.Request, err error) {
			customServeHTTP(logger, te, w, r, err)
		}
	}

	// Send all 404 to the handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Error(http.StatusNotFound, w, r)
	})

	// Set up the router.
	mux.SetServeHTTP(serveHTTP)
	mux.SetNotFound(notFound)
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
	Message() string
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code     int
	Err      error
	Friendly string
}

// Error returns the error.
func (se StatusError) Error() string {
	if se.Err != nil {
		return se.Err.Error()
	}

	return ""
}

// Status returns a HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Message returns a optional user friendly error message.
func (se StatusError) Message() string {
	return se.Friendly
}

package grpcp

import (
	"bytes"
	"context"
	"net/http"

	"github.com/ambientkit/ambient"
)

// HandlerImpl .
type HandlerImpl struct {
	Log ambient.Logger
	Map map[string]func(http.ResponseWriter, *http.Request) error
}

// Handle .
func (d *HandlerImpl) Handle(requestID string, method string, path string, headers http.Header, body []byte) (int, string, string) {
	d.Log.Warn("grpc-plugin: Handle start: %v %v | Routes: %v | %v", method, path, len(d.Map), requestID)

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	newContext := context.WithValue(req.Context(), ambientRequestID, requestID)
	req = req.WithContext(newContext)
	req.Header = headers
	w := NewResponseWriter()

	fn, found := d.Map[pathkey(method, path)]
	if !found {
		return http.StatusNotFound, "", ""
	}

	err := fn(w, req)

	statusCode := 200
	errText := ""
	if err != nil {
		switch e := err.(type) {
		case ambient.Error:
			statusCode = e.Status()
		default:
			statusCode = http.StatusInternalServerError
		}
		errText = err.Error()
		if len(errText) == 0 {
			errText = http.StatusText(statusCode)
		}
	}

	d.Log.Warn("grpc-plugin: Handle end: %v Output: \"%v\"", statusCode, w.Output())

	return statusCode, errText, w.Output()
}

package grpcp

import (
	"bytes"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/requestuuid"
)

// HandlerImpl .
type HandlerImpl struct {
	Log ambient.Logger
	Map map[string]func(http.ResponseWriter, *http.Request) error
}

// Handle .
func (d *HandlerImpl) Handle(requestid string, method string, path string, fullPath string, headers http.Header, body []byte) (int, string, string, http.Header) {
	// d.Log.Warn("grpc-plugin: Handle start: %v %v | Routes: %v | %v", method, path, len(d.Map), requestid)

	req, _ := http.NewRequest(method, fullPath, bytes.NewBuffer(body))
	req = requestuuid.Set(req, requestid)
	req.Header = headers
	w := NewResponseWriter()

	fn, found := d.Map[pathkey(method, path)]
	if !found {
		return http.StatusNotFound, "", "", nil
	}

	err := fn(w, req)

	statusCode := 200
	if w.StatusCode() != 200 {
		statusCode = w.StatusCode()
	}
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

	//d.Log.Warn("grpc-plugin: Handle end: %v Output: \"%v\"", statusCode, w.Output())

	return statusCode, errText, w.Output(), w.Header()
}

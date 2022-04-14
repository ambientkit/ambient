package grpcp

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/grpcsafe"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	"golang.org/x/net/context"
)

// GRPCHandlerPlugin is the gRPC server that GRPCClient talks to.
type GRPCHandlerPlugin struct {
	Log         ambient.Logger
	PluginState *grpcsafe.PluginState
}

// Handle .
func (m *GRPCHandlerPlugin) Handle(ctx context.Context, req *protodef.HandleRequest) (resp *protodef.HandleResponse, err error) {
	headers := http.Header{}
	err = ProtobufStructToObject(req.Headers, &headers)
	if err != nil {
		m.Log.Error("error getting headers: %v", err.Error())
	}

	//status, errText, response, rawHeaders := m.Process(req.Requestid, req.Method, req.Path, req.Fullpath, headers, req.Body)
	status, errText, response, rawHeaders := m.Process(ctx, req, headers)

	outHeaders, err := ObjectToProtobufStruct(rawHeaders)
	if err != nil {
		m.Log.Error("error converting headers: %v", err.Error(), rawHeaders)
		//return &protodef.MiddlewareResponse{}, err
	}

	return &protodef.HandleResponse{
		Status:   uint32(status),
		Error:    errText,
		Response: response,
		Headers:  outHeaders,
	}, err
}

// Process handler.
func (m *GRPCHandlerPlugin) Process(gctx context.Context, hr *protodef.HandleRequest, headers http.Header) (int, string, string, http.Header) {
	// d.Log.Warn("Handle start: %v %v | Routes: %v | %v", method, path, len(d.Map), requestid)

	req := httptest.NewRequest(hr.Method, hr.Fullpath, bytes.NewReader(hr.Body))
	req = req.WithContext(gctx)
	req = requestuuid.Set(req, hr.Requestid)
	req.Header = headers

	// Get the context if saved from middleware.
	ctx, ok := m.PluginState.Context(hr.Requestid)
	if ok {
		req = req.WithContext(ctx)
	}

	w := NewResponseWriter()

	fn, found := m.PluginState.HandleMap(hr.Method, hr.Path)
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

	//d.Log.Warn("Handle end: %v Output: \"%v\"", statusCode, w.Output())

	return statusCode, errText, w.Output(), w.Header()
}

package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// Handler .
type Handler interface {
	Handle(requestID string, method string, path string, fullPath string, headers http.Header, body []byte) (status int, errorText string, response string, outHeaders http.Header)
}

// GRPCHandlerPlugin is the gRPC server that GRPCClient talks to.
type GRPCHandlerPlugin struct {
	Impl Handler
	Log  ambient.Logger
}

// Handle .
func (m *GRPCHandlerPlugin) Handle(ctx context.Context, req *protodef.HandleRequest) (resp *protodef.HandleResponse, err error) {
	headers := http.Header{}
	err = ProtobufStructToObject(req.Headers, &headers)
	if err != nil {
		m.Log.Error("error getting headers: %v", err.Error())
	}

	status, errText, response, rawHeaders := m.Impl.Handle(req.Requestid, req.Method, req.Path, req.Fullpath, headers, req.Body)

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

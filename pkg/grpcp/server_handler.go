package grpcp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCHandlerServer .
type GRPCHandlerServer struct {
	client protodef.HandlerClient
}

// Handle sends the request information from the server to the plugin.
func (l *GRPCHandlerServer) Handle(method string, path string, r *http.Request, requestID string) (
	status int, errText string, response string, headers http.Header, err error) {
	ctx := r.Context()

	sm, err := ObjectToProtobufStruct(r.Header)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", headers, err
	}

	body := bytes.NewBuffer(nil)
	_, err = io.Copy(body, r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", headers, err
	}
	// Restore body.
	r.Body = ioutil.NopCloser(body)

	resp, err := l.client.Handle(ctx, &protodef.HandleRequest{
		Requestid: requestID,
		Method:    method,
		Path:      path,
		Fullpath:  r.URL.RequestURI(),
		Headers:   sm,
		Body:      body.Bytes(),
	})
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", headers, err
	}

	err = ProtobufStructToObject(resp.Headers, &headers)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", headers, err
	}

	return int(resp.Status), resp.Error, resp.Response, headers, err
}

package grpcp

import (
	"io/ioutil"
	"net/http"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCHandlerServer .
type GRPCHandlerServer struct {
	client protodef.HandlerClient
}

// Handle sends the request information from the server to the plugin.
func (l *GRPCHandlerServer) Handle(method string, path string, r *http.Request, requestID string) (
	status int, errText string, response string, err error) {
	ctx := context.Background()

	sm, err := ObjectToProtobufStruct(r.Header)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", err
	}

	resp, err := l.client.Handle(ctx, &protodef.HandleRequest{
		Requestid: requestID,
		Method:    method,
		Path:      path,
		Headers:   sm,
		Body:      body,
	})
	if err != nil {
		return http.StatusInternalServerError, err.Error(), "", err
	}

	return int(resp.Status), resp.Error, resp.Response, err
}

package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/grpcsafe"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	"golang.org/x/net/context"
)

// GRPCRouterPlugin .
type GRPCRouterPlugin struct {
	client      protodef.RouterClient
	Log         ambient.Logger
	PluginState *grpcsafe.PluginState
}

// Handle request handler.
func (c *GRPCRouterPlugin) Handle(method string, path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	//c.Log.Warn("%v called: %v", method, path)
	c.PluginState.SaveHandleMap(fn, method, path)

	c.client.Handle(context.Background(), &protodef.RouterRequest{
		//Uid:    c.brokerID, // TODO: Remove this from the definition.
		Method: method,
		Path:   path,
	})
}

// Get request handler.
func (c *GRPCRouterPlugin) Get(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodGet, path, fn)
}

// Post request handler.
func (c *GRPCRouterPlugin) Post(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodPost, path, fn)
}

// Patch request handler.
func (c *GRPCRouterPlugin) Patch(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodPatch, path, fn)
}

// Put request handler.
func (c *GRPCRouterPlugin) Put(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodPut, path, fn)
}

// Head request handler.
func (c *GRPCRouterPlugin) Head(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodHead, path, fn)
}

// Options request handler.
func (c *GRPCRouterPlugin) Options(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodOptions, path, fn)
}

// Delete request handler.
func (c *GRPCRouterPlugin) Delete(path string, fn func(http.ResponseWriter, *http.Request) (err error)) {
	c.Handle(http.MethodDelete, path, fn)
}

// Param request handler.
func (c *GRPCRouterPlugin) Param(r *http.Request, name string) string {
	v := requestuuid.Get(r)

	out, _ := c.client.Param(context.Background(), &protodef.RouterParamRequest{
		Name:      name,
		Requestid: v,
	})
	return out.Value
}

// StatusError handler.
func (c *GRPCRouterPlugin) StatusError(status int, err error) error {
	return ambient.StatusError{Code: status, Err: err}
}

// Error handler.
func (c *GRPCRouterPlugin) Error(status int, w http.ResponseWriter, r *http.Request) {
	v := requestuuid.Get(r)

	c.Log.Warn("Error called: %v %v", v, status)

	_, err := c.client.Error(context.Background(), &protodef.RouterErrorRequest{
		Status:    uint32(status),
		Requestid: v,
	})
	if err != nil {
		c.Log.Error("error from Error call: %v", err.Error())
	}

}

// Wrap for http.HandlerFunc.
func (c *GRPCRouterPlugin) Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (err error) {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		handler.ServeHTTP(w, r)
		return nil
	}
}

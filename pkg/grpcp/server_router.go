package grpcp

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCAddRouterServer is the gRPC server that GRPCClient talks to.
type GRPCAddRouterServer struct {
	Impl          ambient.Router
	Log           ambient.Logger
	broker        *plugin.GRPCBroker
	conn          *grpc.ClientConn
	HandlerClient *GRPCHandlerServer
	reqmap        *RequestMap
}

// Handle request handler.
func (m *GRPCAddRouterServer) Handle(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	//m.Log.Warn("grpc-server: GET called: %v", req.Path)

	m.Impl.Handle(req.Method, req.Path, func(w http.ResponseWriter, r *http.Request) error {
		//m.Log.Warn("grpc-server: %v func called: %v", req.Method, req.Path)

		uuid := requestuuid.Get(r)
		m.reqmap.Save(uuid, &HTTPContainer{
			Request:  r,
			Response: w,
			FuncMap:  make(template.FuncMap),
		})

		status, errText, response, headers, err := m.HandlerClient.Handle(req.Method, req.Path, r, uuid)
		m.reqmap.Delete(uuid)
		if err != nil {
			m.Log.Error("grpc-server: %v func error: %v", req.Method, err.Error())
			return err
		}

		// TODO: Need to add in headers.

		if status >= 400 && len(response) == 0 {
			return ambient.StatusError{Code: status, Err: errors.New(errText)}
		} else if status == http.StatusFound || status == http.StatusMovedPermanently {
			loc := headers.Get("Location")
			if len(loc) > 0 {
				http.Redirect(w, r, loc, status)
				return nil
			}
			http.Redirect(w, r, response, status)
			return nil
		} else if len(errText) > 0 {
			return errors.New(errText)
		}

		// Only write to response if there is content. The response could have
		// already been handled by other functions like Error().
		if len(response) > 0 {
			w.WriteHeader(status)
			fmt.Fprint(w, response)
		}

		return nil
	})
	return &protodef.Empty{}, err
}

// Get request handler.
func (m *GRPCAddRouterServer) Get(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Post request handler.
func (m *GRPCAddRouterServer) Post(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Patch request handler.
func (m *GRPCAddRouterServer) Patch(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Put request handler.
func (m *GRPCAddRouterServer) Put(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Head request handler.
func (m *GRPCAddRouterServer) Head(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Options request handler.
func (m *GRPCAddRouterServer) Options(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Delete request handler.
func (m *GRPCAddRouterServer) Delete(ctx context.Context, req *protodef.RouterRequest) (resp *protodef.Empty, err error) {
	return m.Handle(ctx, req)
}

// Param returns the named parameters.
func (m *GRPCAddRouterServer) Param(ctx context.Context, req *protodef.RouterParamRequest) (resp *protodef.RouterParamResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.RouterParamResponse{
			Value: "",
		}, nil
	}

	return &protodef.RouterParamResponse{
		Value: m.Impl.Param(c.Request, req.Name),
	}, nil
}

// Delete request handler.
func (m *GRPCAddRouterServer) Error(ctx context.Context, req *protodef.RouterErrorRequest) (resp *protodef.Empty, err error) {
	m.Log.Warn("grpc-server: Error called: %v %v", req.Requestid, req.Status)

	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, nil
	}

	m.Impl.Error(int(req.Status), c.Response, c.Request)
	return &protodef.Empty{}, nil
}

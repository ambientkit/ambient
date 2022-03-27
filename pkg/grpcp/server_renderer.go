package grpcp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/avfs"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCRendererServer is the gRPC server that GRPCClient talks to.
type GRPCRendererServer struct {
	Log              ambient.Logger
	Impl             ambient.Renderer
	reqmap           *RequestMap
	FuncMapperClient *GRPCFuncMapperServer
}

// Page handler.
func (m *GRPCRendererServer) Page(ctx context.Context, req *protodef.RendererPageRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	// Build a file system.
	efs := avfs.NewFS()
	for _, v := range req.Files {
		efs.AddFile(v.Name, v.Body)
	}

	vars := make(map[string]interface{})
	err = ProtobufStructToObject(req.Vars, &vars)
	if err != nil {
		return &protodef.Empty{}, fmt.Errorf("grpc-server: error on Page object conversion: %v", err.Error())
	}

	err = m.Impl.Page(c.Response, c.Request, efs, req.Templatename, func(*http.Request) template.FuncMap {
		fm := template.FuncMap{}
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			fm[v] = func(args ...interface{}) (interface{}, error) {
				//m.Log.Error("FUNC ARGS: %v | %#v", v, args)
				val, err := m.FuncMapperClient.Do(c.Request, req.Requestid, v, args, false)
				return val, err
			}
		}
		return fm
	}, vars)

	return &protodef.Empty{}, err
}

// PageContent handler.
func (m *GRPCRendererServer) PageContent(ctx context.Context, req *protodef.RendererPageContentRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	vars := make(map[string]interface{})
	err = ProtobufStructToObject(req.Vars, &vars)
	if err != nil {
		return &protodef.Empty{}, fmt.Errorf("grpc-server: error on PageContent object conversion: %v", err.Error())
	}

	err = m.Impl.PageContent(c.Response, c.Request, req.Content, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(c.Request, req.Requestid, v, args, false)
				return val, err
			}
		}
		return c.FuncMap
	}, vars)

	return &protodef.Empty{}, err
}

// Post handler.
func (m *GRPCRendererServer) Post(ctx context.Context, req *protodef.RendererPostRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	// Build a file system.
	efs := avfs.NewFS()
	for _, v := range req.Files {
		efs.AddFile(v.Name, v.Body)
	}

	vars := make(map[string]interface{})
	err = ProtobufStructToObject(req.Vars, &vars)
	if err != nil {
		return &protodef.Empty{}, fmt.Errorf("grpc-server: error on Post object conversion: %v", err.Error())
	}

	err = m.Impl.Post(c.Response, c.Request, efs, req.Templatename, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(c.Request, req.Requestid, v, args, false)
				return val, err
			}
		}
		return c.FuncMap
	}, vars)

	return &protodef.Empty{}, err
}

// PostContent handler.
func (m *GRPCRendererServer) PostContent(ctx context.Context, req *protodef.RendererPostContentRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	vars := make(map[string]interface{})
	err = ProtobufStructToObject(req.Vars, &vars)
	if err != nil {
		return &protodef.Empty{}, fmt.Errorf("grpc-server: error on PostContent object conversion: %v", err.Error())
	}

	err = m.Impl.PostContent(c.Response, c.Request, req.Content, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(c.Request, req.Requestid, v, args, false)
				return val, err
			}
		}
		return c.FuncMap
	}, vars)

	return &protodef.Empty{}, err
}

// Error handler.
func (m *GRPCRendererServer) Error(ctx context.Context, req *protodef.RendererErrorRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	vars := make(map[string]interface{})
	err = ProtobufStructToObject(req.Vars, &vars)
	if err != nil {
		return &protodef.Empty{}, fmt.Errorf("grpc-server: error on Error object conversion: %v", err.Error())
	}

	err = m.Impl.Error(c.Response, c.Request, req.Content, int(req.Statuscode), func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(c.Request, req.Requestid, v, args, false)
				return val, err
			}
		}
		return c.FuncMap
	}, vars)

	return &protodef.Empty{}, err
}

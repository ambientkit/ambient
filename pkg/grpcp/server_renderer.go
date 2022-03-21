package grpcp

import (
	"embed"
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

	//efs := embed.FS{}
	//efs.

	efs := avfs.NewFS()
	efs.AddFile("template/content/hello.tmpl", []byte("woah"))
	//req.Files

	_, err = efs.Open("template/content/hello.tmpl")
	if err != nil {
		return nil, err
	}

	//fff := new(fs.FS)
	//fff.

	err = m.Impl.Page(c.Response, c.Request, efs, req.Templatename, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(req.Requestid, v, args)
				return val, err
			}
		}
		return c.FuncMap
	}, ProtobufStructToMap(req.Vars))

	return &protodef.Empty{}, err
}

// PageContent handler.
func (m *GRPCRendererServer) PageContent(ctx context.Context, req *protodef.RendererPageContentRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.PageContent(c.Response, c.Request, req.Content, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(req.Requestid, v, args)
				return val, err
			}
		}
		return c.FuncMap
	}, ProtobufStructToMap(req.Vars))

	return &protodef.Empty{}, err
}

// Post handler.
func (m *GRPCRendererServer) Post(ctx context.Context, req *protodef.RendererPostRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.Post(c.Response, c.Request, embed.FS{}, req.Templatename, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(req.Requestid, v, args)
				return val, err
			}
		}
		return c.FuncMap
	}, ProtobufStructToMap(req.Vars))

	return &protodef.Empty{}, err
}

// PostContent handler.
func (m *GRPCRendererServer) PostContent(ctx context.Context, req *protodef.RendererPostContentRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.PostContent(c.Response, c.Request, req.Content, func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(req.Requestid, v, args)
				return val, err
			}
		}
		return c.FuncMap
	}, ProtobufStructToMap(req.Vars))

	return &protodef.Empty{}, err
}

// Error handler.
func (m *GRPCRendererServer) Error(ctx context.Context, req *protodef.RendererErrorRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.Error(c.Response, c.Request, req.Content, int(req.Statuscode), func(*http.Request) template.FuncMap {
		for _, rawV := range req.Keys {
			// Prevent race conditions.
			v := rawV
			c.FuncMap[v] = func(args ...interface{}) (interface{}, error) {
				val, err := m.FuncMapperClient.Do(req.Requestid, v, args)
				return val, err
			}
		}
		return c.FuncMap
	}, ProtobufStructToMap(req.Vars))

	return &protodef.Empty{}, err
}

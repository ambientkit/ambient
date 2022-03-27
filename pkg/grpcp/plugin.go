package grpcp

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCPlugin is the plugin side implementation.
type GRPCPlugin struct {
	Impl             ambient.MiddlewarePlugin
	broker           *plugin.GRPCBroker
	toolkit          *ambient.Toolkit
	conn             *grpc.ClientConn
	server           *grpc.Server
	funcMapperClient *GRPCFuncMapperServer
	reqMap           map[string]func(http.ResponseWriter, *http.Request) error
	funcMap          map[string]*FMContainer
}

// PluginName handler.
func (m *GRPCPlugin) PluginName(ctx context.Context, req *protodef.Empty) (*protodef.PluginNameResponse, error) {
	name := m.Impl.PluginName()
	return &protodef.PluginNameResponse{Name: name}, nil
}

// PluginVersion handler.
func (m *GRPCPlugin) PluginVersion(ctx context.Context, req *protodef.Empty) (*protodef.PluginVersionResponse, error) {
	version := m.Impl.PluginVersion()
	return &protodef.PluginVersionResponse{Version: version}, nil
}

// GrantRequests handler.
func (m *GRPCPlugin) GrantRequests(ctx context.Context, req *protodef.Empty) (*protodef.GrantRequestsResponse, error) {
	requests := m.Impl.GrantRequests()
	arr := make([]*protodef.GrantRequest, 0)

	for _, v := range requests {
		arr = append(arr, &protodef.GrantRequest{
			Description: v.Description,
			Grant:       string(v.Grant),
		})
	}

	return &protodef.GrantRequestsResponse{Grantrequest: arr}, nil
}

// Enable handler.
func (m *GRPCPlugin) Enable(ctx context.Context, req *protodef.Toolkit) (*protodef.EnableResponse, error) {
	var err error
	m.conn, err = m.broker.Dial(req.Uid)
	if err != nil {
		return nil, err
	}

	logger := &GRPCLoggerPlugin{
		client: protodef.NewLoggerClient(m.conn),
	}

	m.funcMapperClient = &GRPCFuncMapperServer{
		client: protodef.NewFuncMapperClient(m.conn),
	}

	m.reqMap = make(map[string]func(http.ResponseWriter, *http.Request) error)
	m.funcMap = make(map[string]*FMContainer)

	m.toolkit = &ambient.Toolkit{
		Log: logger,
		Mux: &GRPCRouterPlugin{
			client: protodef.NewRouterClient(m.conn),
			Log:    logger,
			Map:    m.reqMap,
		},
		Site: &GRPCSitePlugin{
			client: protodef.NewSiteClient(m.conn),
			Log:    logger,
		},
		Render: &GRPCRendererPlugin{
			client: protodef.NewRendererClient(m.conn),
			Log:    logger,
			Map:    m.funcMap,
		},
	}

	m.toolkit.Log.Debug("grpc-plugin: Enabled() called")

	err = m.Impl.Enable(m.toolkit)

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		m.server = grpc.NewServer(opts...)
		protodef.RegisterFuncMapperServer(m.server, &GRPCFuncMapperPlugin{
			Impl: &FuncMapperImpl{
				Log:  m.toolkit.Log,
				Map:  m.funcMap,
				Impl: m.Impl,
			},
			Log: m.toolkit.Log,
		})
		protodef.RegisterHandlerServer(m.server, &GRPCHandlerPlugin{
			Log: m.toolkit.Log,
			Impl: &HandlerImpl{
				Log: m.toolkit.Log,
				Map: m.reqMap,
			},
		})

		return m.server
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	return &protodef.EnableResponse{
		Uid: brokerID,
	}, err
}

// Disable handler.
func (m *GRPCPlugin) Disable(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	m.toolkit.Log.Debug("grpc-plugin: Disable() called")
	defer m.conn.Close()
	return &protodef.Empty{}, m.Impl.Disable()
}

// Routes handler.
func (m *GRPCPlugin) Routes(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	m.toolkit.Log.Debug("grpc-plugin: Routes() called")
	m.Impl.Routes()
	return &protodef.Empty{}, nil
}

// Assets handler.
func (m *GRPCPlugin) Assets(ctx context.Context, req *protodef.Empty) (*protodef.AssetsResponse, error) {
	settings, embedFS := m.Impl.Assets()

	assets, err := ArrayToProtobufStruct(settings)
	if err != nil {
		return &protodef.AssetsResponse{}, err
	}

	files := make([]*protodef.EmbeddedFile, 0)

	if embedFS == nil {
		return &protodef.AssetsResponse{
			Assets: assets,
			Files:  files,
		}, nil
	}

	err = fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := embedFS.ReadFile(path)
		if err != nil {
			return err
		}

		files = append(files, &protodef.EmbeddedFile{
			Name: path,
			Body: b,
		})
		return nil
	})
	if err != nil {
		return &protodef.AssetsResponse{
			Assets: assets,
			Files:  files,
		}, err
	}

	return &protodef.AssetsResponse{
		Assets: assets,
		Files:  files,
	}, nil
}

// Settings handler.
func (m *GRPCPlugin) Settings(ctx context.Context, req *protodef.Empty) (*protodef.SettingsResponse, error) {
	settings := m.Impl.Settings()

	arr := make([]*protodef.Setting, 0)
	for _, v := range settings {
		any, err := InterfaceToProtobufAny(v.Default)
		if err != nil {
			m.toolkit.Log.Error("grpc-plugin: error on conversion: %v", err)
		}

		arr = append(arr, &protodef.Setting{
			Name:        v.Name,
			Settingtype: string(v.Type),
			Description: &protodef.SettingDescription{
				Text: v.Description.Text,
				Url:  v.Description.URL,
			},
			Hide:    v.Hide,
			Default: any,
		})
	}

	return &protodef.SettingsResponse{
		Settings: arr,
	}, nil
}

// FuncMap handler.
func (m *GRPCPlugin) FuncMap(ctx context.Context, req *protodef.Empty) (*protodef.FuncMapResponse, error) {
	//m.toolkit.Log.Error("grpc-plugin: FuncMap called.")
	fn := m.Impl.FuncMap()
	r := httptest.NewRequest("GET", "/", nil)
	if fn == nil {
		return &protodef.FuncMapResponse{}, nil
	}

	fm := fn(r)

	keys := make([]string, 0)
	for k := range fm {
		keys = append(keys, k)
	}

	return &protodef.FuncMapResponse{
		Keys: keys,
	}, nil
}

// Middleware handler.
func (m *GRPCPlugin) Middleware(ctx context.Context, req *protodef.MiddlewareRequest) (*protodef.MiddlewareResponse, error) {
	m.toolkit.Log.Debug("grpc-plugin: Middleware() called")

	// Get the middleware from the plugin.
	arr := m.Impl.Middleware()
	if len(arr) == 0 {
		return &protodef.MiddlewareResponse{}, nil
	}

	headers := http.Header{}
	err := ProtobufStructToObject(req.Headers, &headers)
	if err != nil {
		return &protodef.MiddlewareResponse{}, err
	}

	r := httptest.NewRequest(req.Method, req.Path, bytes.NewBuffer(req.Body))
	r = requestuuid.Set(r, req.Requestid)
	r.Header = headers
	w := NewResponseWriter()

	mux := &MockHandler{
		//Log: m.toolkit.Log,
	}
	var h http.Handler
	h = mux

	for _, mw := range arr {
		//m.toolkit.Log.Warn("Looping for: %v %v %v", req.Method, req.Path, req.Requestid)
		h = mw(h)
	}
	h.ServeHTTP(w, r)

	//m.toolkit.Log.Warn("Plugin out: %v | %v", mux.W, w.Output())

	statusCode := 0
	if w.statusCode != 0 {
		statusCode = w.statusCode
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

	//m.toolkit.Log.Error("Sending Middleware: %v | %v | %v", statusCode, errText, w.Output())

	outHeaders, err := ObjectToProtobufStruct(w.header)
	if err != nil {
		m.toolkit.Log.Error("grpc-plugin: error getting headers: %v", err.Error())
		return &protodef.MiddlewareResponse{}, err
	}

	return &protodef.MiddlewareResponse{
		Status:   uint32(statusCode),
		Error:    errText,
		Response: w.Output(),
		Headers:  outHeaders,
	}, nil
}

// MockHandler is a mock mux.
type MockHandler struct {
	//W http.ResponseWriter
	//R *http.Request
	//Log ambient.Logger
}

// ServeHTTP stores the requests.
func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//h.Log.Warn("Final loop for: %v %v %v", r.Method, r.URL.Path, requestuuid.Get(r))
	//h.W = w
	//h.R = r
}

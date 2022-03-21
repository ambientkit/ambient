package grpcp

import (
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCPlugin is the plugin side implementation.
type GRPCPlugin struct {
	Impl    ambient.Plugin
	broker  *plugin.GRPCBroker
	toolkit *ambient.Toolkit
	conn    *grpc.ClientConn
	server  *grpc.Server
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

	fnMap := make(map[string]func(http.ResponseWriter, *http.Request) error)

	funcMapMap := make(map[string]template.FuncMap)

	m.toolkit = &ambient.Toolkit{
		Log: logger,
		Mux: &GRPCRouterPlugin{
			client: protodef.NewRouterClient(m.conn),
			Log:    logger,
			Map:    fnMap,
		},
		Site: &GRPCSitePlugin{
			client: protodef.NewSiteClient(m.conn),
			Log:    logger,
		},
		Render: &GRPCRendererPlugin{
			client: protodef.NewRendererClient(m.conn),
			Log:    logger,
			Map:    funcMapMap,
		},
	}

	m.toolkit.Log.Debug("grpc-plugin: Enabled() called")

	err = m.Impl.Enable(m.toolkit)

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		m.server = grpc.NewServer(opts...)
		protodef.RegisterFuncMapperServer(m.server, &GRPCFuncMapperPlugin{
			Impl: &FuncMapperImpl{
				Log: m.toolkit.Log,
				Map: funcMapMap,
			},
			Log: m.toolkit.Log,
		})
		protodef.RegisterHandlerServer(m.server, &GRPCHandlerPlugin{
			Log: m.toolkit.Log,
			Impl: &HandlerImpl{
				Log: m.toolkit.Log,
				Map: fnMap,
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
	settings, _ := m.Impl.Assets()

	assets, err := ArrayToProtobufStruct(settings)
	if err != nil {
		return &protodef.AssetsResponse{}, err
	}

	// arr := make([]*protodef.Setting, 0)
	// for _, v := range settings {
	// 	any, err := InterfaceToProtobufAny(v.Default)
	// 	if err != nil {
	// 		m.toolkit.Log.Error("grpc-plugin: error on conversion: %v", err)
	// 	}

	// 	arr = append(arr, &protodef.Setting{
	// 		Name:        v.Name,
	// 		Settingtype: string(v.Type),
	// 		Description: &protodef.SettingDescription{
	// 			Text: v.Description.Text,
	// 			Url:  v.Description.URL,
	// 		},
	// 		Hide:    v.Hide,
	// 		Default: any,
	// 	})
	// }

	return &protodef.AssetsResponse{
		Assets: assets,
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
func (m *GRPCPlugin) FuncMap(ctx context.Context, req *protodef.FuncMapRequest) (*protodef.FuncMapResponse, error) {
	//fm := m.Impl.FuncMap()
	//fm()

	m.toolkit.Log.Error("grpc-plugin: NEED TO IMPLEMENT. FuncMap() hit! Request id:", req.Requestid)

	fmfn := m.Impl.FuncMap()
	fm := fmfn(nil)

	keys := make([]string, 0)
	for k := range fm {
		keys = append(keys, k)
	}

	return &protodef.FuncMapResponse{
		Keys: keys,
	}, nil
}

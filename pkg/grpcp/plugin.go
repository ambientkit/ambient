package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCPlugin is the plugin side implementation.
type GRPCPlugin struct {
	Impl    PluginCore
	broker  *plugin.GRPCBroker
	toolkit *Toolkit
	conn    *grpc.ClientConn
	server  *grpc.Server
}

// PluginName returns the plugin name.
func (m *GRPCPlugin) PluginName(ctx context.Context, req *protodef.Empty) (*protodef.PluginNameResponse, error) {
	name, err := m.Impl.PluginName()
	return &protodef.PluginNameResponse{Name: name}, err
}

// PluginVersion returns the plugin version.
func (m *GRPCPlugin) PluginVersion(ctx context.Context, req *protodef.Empty) (*protodef.PluginVersionResponse, error) {
	version, err := m.Impl.PluginVersion()
	return &protodef.PluginVersionResponse{Version: version}, err
}

// GrantRequests .
// func (m *GRPCPlugin) GrantRequests(ctx context.Context, req *protodef.Empty) (*protodef.GrantRequestsResponse, error) {
// 	requests := m.Impl.GrantRequests()
// 	arr := make([]*proto.GrantRequest, 0)

// 	for _, v := range requests {
// 		arr = append(arr, &proto.GrantRequest{
// 			Description: v.Description,
// 			Grant:       string(v.Grant),
// 		})
// 	}

// 	return &proto.GrantRequestsResponse{GrantRequest: arr}, nil
// }

// Enable .
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

	m.toolkit = &Toolkit{
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
	}

	m.toolkit.Log.Debug("grpc-plugin: Enabled() called")

	err = m.Impl.Enable(m.toolkit)

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		m.server = grpc.NewServer(opts...)
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

// Disable will disable the plugin and close connections.
func (m *GRPCPlugin) Disable(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	m.toolkit.Log.Debug("grpc-plugin: Disable() called")
	defer m.conn.Close()
	return &protodef.Empty{}, m.Impl.Disable()
}

// Routes .
func (m *GRPCPlugin) Routes(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	m.toolkit.Log.Debug("grpc-plugin: Routes() called")
	return &protodef.Empty{}, m.Impl.Routes()
}

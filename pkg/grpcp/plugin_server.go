package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCPlugin is the gRPC server that GRPCClient talks to.
type GRPCPlugin struct {
	Impl    PluginCore
	broker  *plugin.GRPCBroker
	toolkit *Toolkit
	conn    *grpc.ClientConn
	server  *grpc.Server
}

// PluginName .
func (m *GRPCPlugin) PluginName(ctx context.Context, req *protodef.Empty) (*protodef.PluginNameResponse, error) {
	return &protodef.PluginNameResponse{Name: m.Impl.PluginName()}, nil
}

// PluginVersion .
func (m *GRPCPlugin) PluginVersion(ctx context.Context, req *protodef.Empty) (*protodef.PluginVersionResponse, error) {
	return &protodef.PluginVersionResponse{Version: m.Impl.PluginVersion()}, nil
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

// Disable .
func (m *GRPCPlugin) Disable(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	defer m.conn.Close()
	err := m.Impl.Disable()
	return &protodef.Empty{}, err
}

// Routes .
func (m *GRPCPlugin) Routes(ctx context.Context, req *protodef.Empty) (*protodef.Empty, error) {
	m.toolkit.Log.Warn("grpc-plugin: routes called")
	m.Impl.Routes()
	return &protodef.Empty{}, nil
}

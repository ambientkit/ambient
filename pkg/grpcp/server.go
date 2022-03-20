package grpcp

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCServer is the server side implementation.
type GRPCServer struct {
	broker  *plugin.GRPCBroker
	client  protodef.GenericPluginClient
	toolkit *ambient.Toolkit
	conn    *grpc.ClientConn
	server  *grpc.Server
	reqmap  *RequestMap
}

// PluginName handler.
func (m *GRPCServer) PluginName() string {
	resp, err := m.client.PluginName(context.Background(), &protodef.Empty{})
	if err != nil {
		return ""
	}

	return resp.Name
}

// PluginVersion handler.
func (m *GRPCServer) PluginVersion() string {
	resp, err := m.client.PluginVersion(context.Background(), &protodef.Empty{})
	if err != nil {
		return ""
	}

	return resp.Version
}

// Enable handler.
func (m *GRPCServer) Enable(toolkit *ambient.Toolkit) error {
	toolkit.Log.Debug("grpc-server: enabled called")

	m.reqmap = NewRequestMap()
	m.toolkit = toolkit
	loggerServer := &GRPCLoggerServer{Impl: toolkit.Log}
	routerServer := &GRPCAddRouterServer{Impl: toolkit.Mux,
		Log:    toolkit.Log,
		broker: m.broker,
		reqmap: m.reqmap,
	}
	siteServer := &GRPCSiteServer{Impl: toolkit.Site, Log: toolkit.Log, reqmap: m.reqmap}

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		m.server = grpc.NewServer(opts...)
		protodef.RegisterLoggerServer(m.server, loggerServer)
		protodef.RegisterRouterServer(m.server, routerServer)
		protodef.RegisterSiteServer(m.server, siteServer)

		return m.server
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	resp, err := m.client.Enable(context.Background(), &protodef.Toolkit{
		Uid: brokerID,
	})
	if err != nil {
		return err
	}

	m.conn, err = m.broker.Dial(resp.Uid)
	if err != nil {
		return err
	}

	routerServer.HandlerClient = &GRPCHandlerServer{
		client: protodef.NewHandlerClient(m.conn),
	}

	return nil
}

// Disable handler.
func (m *GRPCServer) Disable() error {
	if m.server != nil {
		_, _ = m.client.Disable(context.Background(), &protodef.Empty{})
		m.server.Stop()
		m.server = nil
	}
	return nil
}

// Routes handler.
func (m *GRPCServer) Routes() {
	if m.server == nil || m.toolkit == nil || m.toolkit.Log == nil {
		return //fmt.Errorf("grpc-server: plugin is disabled")
	}

	m.toolkit.Log.Warn("grpc-server: routes called")

	_, err := m.client.Routes(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("grpc-server: error calling routes: %v", err)
	}

	m.toolkit.Log.Warn("grpc-server: routes called END")

	//return err
}

// Assets handler.
func (m *GRPCServer) Assets() ([]ambient.Asset, *embed.FS) {
	// TODO: Implement
	return nil, nil
}

// Settings handler.
func (m *GRPCServer) Settings() []ambient.Setting {
	resp, err := m.client.Settings(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("grpc-server: error calling Settings: %v", err)
	}

	arr := make([]ambient.Setting, 0)

	for _, v := range resp.Settings {
		var i interface{}
		err = ProtobufAnyToInterface(v.Default, &i)
		if err != nil {
			m.toolkit.Log.Error("grpc-server: error calling Settings: %v", err)
		}

		arr = append(arr, ambient.Setting{
			Name: v.Name,
			Type: ambient.SettingType(v.Settingtype),
			Description: ambient.SettingDescription{
				Text: v.Description.Text,
				URL:  v.Description.Url,
			},
			Hide:    v.Hide,
			Default: i,
		})
	}

	if err != nil {
		m.toolkit.Log.Error("grpc-server: error calling Settings conversion: %v", err)
	}
	return arr
}

// GrantRequests handler.
func (m *GRPCServer) GrantRequests() []ambient.GrantRequest {
	resp, err := m.client.GrantRequests(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("grpc-server: error calling GrantRequests: %v", err)
	}

	arr := make([]ambient.GrantRequest, 0)
	for _, v := range resp.GrantRequest {
		arr = append(arr, ambient.GrantRequest{
			Grant:       ambient.Grant(v.Grant),
			Description: v.Description,
		})
	}

	return arr
}

// FuncMap handler.
func (m *GRPCServer) FuncMap() func(r *http.Request) template.FuncMap {
	// TODO: Implement
	return nil
}

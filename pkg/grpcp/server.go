package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCServer is the server side implementation.
type GRPCServer struct {
	broker  *plugin.GRPCBroker
	client  protodef.GenericPluginClient
	toolkit *Toolkit
	conn    *grpc.ClientConn
	server  *grpc.Server
	reqmap  *RequestMap
}

// PluginName .
func (m *GRPCServer) PluginName() (string, error) {
	resp, err := m.client.PluginName(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Name, nil
}

// PluginVersion .
func (m *GRPCServer) PluginVersion() (string, error) {
	resp, err := m.client.PluginVersion(context.Background(), &protodef.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

// Enable .
func (m *GRPCServer) Enable(toolkit *Toolkit) error {
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

// Disable .
func (m *GRPCServer) Disable() error {
	if m.server != nil {
		_, _ = m.client.Disable(context.Background(), &protodef.Empty{})
		m.server.Stop()
		m.server = nil
	}
	return nil
}

// Routes .
func (m *GRPCServer) Routes() error {
	if m.server == nil || m.toolkit == nil || m.toolkit.Log == nil {
		return fmt.Errorf("grpc-server: plugin is disabled")
	}

	m.toolkit.Log.Warn("grpc-server: routes called")

	_, err := m.client.Routes(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("grpc-server: error calling routes: %v", err)
	}

	m.toolkit.Log.Warn("grpc-server: routes called END")

	return err
}

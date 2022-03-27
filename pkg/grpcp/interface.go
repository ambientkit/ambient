// Package grpcp contains shared data between the host and plugins.
package grpcp

import (
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "AMBIENT_PLUGIN",
	MagicCookieValue: "v1.0",
}

// GenericPlugin is the implementation of plugin. Plugin so we can serve/consume
// this. We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type GenericPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ambient.MiddlewarePlugin
}

// GRPCServer .
func (p *GenericPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	protodef.RegisterGenericPluginServer(s, &GRPCPlugin{
		Impl:       p.Impl,
		broker:     broker,
		contextMap: make(map[string]context.Context),
		reqMap:     make(map[string]func(http.ResponseWriter, *http.Request) error),
		funcMap:    make(map[string]*FMContainer),
	})
	return nil
}

// GRPCClient .
func (p *GenericPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCServer{
		client: protodef.NewGenericPluginClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &GenericPlugin{}

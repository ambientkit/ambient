package grpcp

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCLoggerServer is the gRPC server that GRPCClient talks to.
type GRPCLoggerServer struct {
	Impl ambient.Logger
}

// Debug handler.
func (m *GRPCLoggerServer) Debug(ctx context.Context, req *protodef.LogFormat) (resp *protodef.Empty, err error) {
	//m.Impl.For(ctx).Debug(req.Format)
	m.Impl.Debug(req.Format)
	return &protodef.Empty{}, err
}

// Info handler.
func (m *GRPCLoggerServer) Info(ctx context.Context, req *protodef.LogFormat) (resp *protodef.Empty, err error) {
	//m.Impl.For(ctx).Info(req.Format)
	m.Impl.Info(req.Format)
	return &protodef.Empty{}, err
}

// Warn handler.
func (m *GRPCLoggerServer) Warn(ctx context.Context, req *protodef.LogFormat) (resp *protodef.Empty, err error) {
	m.Impl.Warn(req.Format)
	return &protodef.Empty{}, err
}

// Error handler.
func (m *GRPCLoggerServer) Error(ctx context.Context, req *protodef.LogFormat) (resp *protodef.Empty, err error) {
	m.Impl.Error(req.Format)
	return &protodef.Empty{}, err
}

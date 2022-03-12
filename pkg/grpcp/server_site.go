package grpcp

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
)

// GRPCSiteServer is the gRPC server that GRPCClient talks to.
type GRPCSiteServer struct {
	Impl   SecureSite
	Log    ambient.Logger
	reqmap *RequestMap
}

// UserLogin .
func (m *GRPCSiteServer) UserLogin(ctx context.Context, req *protodef.SiteUserLoginRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return
	}
	err = m.Impl.UserLogin(c.Request, req.Username)
	return &protodef.Empty{}, err
}

// AuthenticatedUser .
func (m *GRPCSiteServer) AuthenticatedUser(ctx context.Context, req *protodef.SiteAuthenticatedUserRequest) (resp *protodef.SiteAuthenticatedUserResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return
	}
	username, err := m.Impl.AuthenticatedUser(c.Request)
	return &protodef.SiteAuthenticatedUserResponse{
		Username: username,
	}, err
}

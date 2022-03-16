package grpcp

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPCSiteServer is the server side implementation of secure site.
type GRPCSiteServer struct {
	Impl   SecureSite
	Log    ambient.Logger
	reqmap *RequestMap
}

// Load handler.
func (m *GRPCSiteServer) Load(ctx context.Context, req *protodef.Empty) (resp *protodef.Empty, err error) {
	err = m.Impl.Load()
	return &protodef.Empty{}, err
}

// Authorized handler.
func (m *GRPCSiteServer) Authorized(ctx context.Context, req *protodef.SiteAuthorizedRequest) (resp *protodef.SiteAuthorizedResponse, err error) {
	authorized := m.Impl.Authorized(ambient.Grant(req.Grant))
	return &protodef.SiteAuthorizedResponse{
		Authorized: authorized,
	}, err
}

// NeighborPluginGrantList handler.
func (m *GRPCSiteServer) NeighborPluginGrantList(ctx context.Context, req *protodef.SiteNeighborPluginGrantListRequest) (
	resp *protodef.SiteNeighborPluginGrantListResponse, err error) {
	gr, err := m.Impl.NeighborPluginGrantList(req.Pluginname)
	if err != nil {
		return &protodef.SiteNeighborPluginGrantListResponse{
			Grants: []*protodef.GrantRequest{},
		}, err
	}

	arr := make([]*protodef.GrantRequest, 0)
	for _, v := range gr {
		arr = append(arr, &protodef.GrantRequest{
			Grant:       string(v.Grant),
			Description: v.Description,
		})
	}

	return &protodef.SiteNeighborPluginGrantListResponse{
		Grants: arr,
	}, err
}

// NeighborPluginGrants handler.
func (m *GRPCSiteServer) NeighborPluginGrants(ctx context.Context, req *protodef.SiteNeighborPluginGrantsRequest) (
	resp *protodef.SiteNeighborPluginGrantsResponse, err error) {
	gr, err := m.Impl.NeighborPluginGrants(req.Pluginname)
	if err != nil {
		return &protodef.SiteNeighborPluginGrantsResponse{
			Grants: &structpb.Struct{},
		}, err
	}

	arr, err := GrantBoolMapToProtobufStruct(gr)
	return &protodef.SiteNeighborPluginGrantsResponse{
		Grants: arr,
	}, err
}

/////////////////////////////////////////////////////

// UserLogin handler.
func (m *GRPCSiteServer) UserLogin(ctx context.Context, req *protodef.SiteUserLoginRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return
	}
	err = m.Impl.UserLogin(c.Request, req.Username)
	return &protodef.Empty{}, err
}

// AuthenticatedUser handler.
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

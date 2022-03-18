package grpcp

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPCSiteServer is the server side implementation of secure site.
type GRPCSiteServer struct {
	Impl   ambient.SecureSite
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

	arr, err := ObjectToProtobufStruct(gr)
	return &protodef.SiteNeighborPluginGrantsResponse{
		Grants: arr,
	}, err
}

// NeighborPluginGranted handler.
func (m *GRPCSiteServer) NeighborPluginGranted(ctx context.Context, req *protodef.SiteNeighborPluginGrantedRequest) (
	resp *protodef.SiteNeighborPluginGrantedResponse, err error) {
	granted, err := m.Impl.NeighborPluginGranted(req.Pluginname, ambient.Grant(req.Grant))
	if err != nil {
		return &protodef.SiteNeighborPluginGrantedResponse{}, err
	}

	return &protodef.SiteNeighborPluginGrantedResponse{
		Granted: granted,
	}, err
}

// NeighborPluginRequestedGrant handler.
func (m *GRPCSiteServer) NeighborPluginRequestedGrant(ctx context.Context, req *protodef.SiteNeighborPluginRequestedGrantRequest) (
	resp *protodef.SiteNeighborPluginRequestedGrantResponse, err error) {
	granted, err := m.Impl.NeighborPluginRequestedGrant(req.Pluginname, ambient.Grant(req.Grant))
	if err != nil {
		return &protodef.SiteNeighborPluginRequestedGrantResponse{}, err
	}

	return &protodef.SiteNeighborPluginRequestedGrantResponse{
		Granted: granted,
	}, err
}

// SetNeighborPluginGrant handler.
func (m *GRPCSiteServer) SetNeighborPluginGrant(ctx context.Context, req *protodef.SiteSetNeighborPluginGrantRequest) (
	resp *protodef.Empty, err error) {
	err = m.Impl.SetNeighborPluginGrant(req.Pluginname, ambient.Grant(req.Grant), req.Granted)
	return &protodef.Empty{}, err
}

// Plugins handler.
func (m *GRPCSiteServer) Plugins(ctx context.Context, req *protodef.Empty) (
	resp *protodef.SitePluginsResponse, err error) {
	pd, err := m.Impl.Plugins()
	if err != nil {
		return &protodef.SitePluginsResponse{
			Plugindata: &structpb.Struct{},
		}, err
	}

	arr, err := ObjectToProtobufStruct(pd)
	return &protodef.SitePluginsResponse{
		Plugindata: arr,
	}, err
}

// PluginNames handler.
func (m *GRPCSiteServer) PluginNames(ctx context.Context, req *protodef.Empty) (
	resp *protodef.SitePluginNamesResponse, err error) {
	names, err := m.Impl.PluginNames()
	if err != nil {
		return &protodef.SitePluginNamesResponse{
			Names: make([]string, 0),
		}, err
	}

	return &protodef.SitePluginNamesResponse{
		Names: names,
	}, err
}

// DeletePlugin handler.
func (m *GRPCSiteServer) DeletePlugin(ctx context.Context, req *protodef.SiteDeletePluginRequest) (
	resp *protodef.Empty, err error) {
	err = m.Impl.DeletePlugin(req.Name)
	return &protodef.Empty{}, err
}

// EnablePlugin handler.
func (m *GRPCSiteServer) EnablePlugin(ctx context.Context, req *protodef.SiteEnablePluginRequest) (
	resp *protodef.Empty, err error) {
	err = m.Impl.EnablePlugin(req.Name, req.Load)
	return &protodef.Empty{}, err
}

// LoadAllPluginPages handler.
func (m *GRPCSiteServer) LoadAllPluginPages(ctx context.Context, req *protodef.Empty) (
	resp *protodef.Empty, err error) {
	err = m.Impl.LoadAllPluginPages()
	return &protodef.Empty{}, err
}

// DisablePlugin handler.
func (m *GRPCSiteServer) DisablePlugin(ctx context.Context, req *protodef.SiteDisablePluginRequest) (
	resp *protodef.Empty, err error) {
	err = m.Impl.DisablePlugin(req.Name, req.Unload)
	return &protodef.Empty{}, err
}

// SavePost handler.
func (m *GRPCSiteServer) SavePost(ctx context.Context, req *protodef.SiteSavePostRequest) (
	resp *protodef.Empty, err error) {
	post := ambient.Post{}
	err = ProtobufStructToObject(req.Post, &post)
	if err != nil {
		return &protodef.Empty{}, err
	}
	err = m.Impl.SavePost(req.Id, post)
	return &protodef.Empty{}, err
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

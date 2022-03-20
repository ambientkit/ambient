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

// PostsAndPages handler.
func (m *GRPCSiteServer) PostsAndPages(ctx context.Context, req *protodef.SitePostsAndPagesRequest) (
	resp *protodef.SitePostsAndPagesResponse, err error) {
	post, err := m.Impl.PostsAndPages(req.Onlypublished)
	if err != nil {
		return &protodef.SitePostsAndPagesResponse{}, err
	}

	p, err := ArrayToProtobufStruct(post)
	return &protodef.SitePostsAndPagesResponse{
		Postwithidlist: p,
	}, err
}

// PublishedPosts handler.
func (m *GRPCSiteServer) PublishedPosts(ctx context.Context, req *protodef.Empty) (
	resp *protodef.SitePublishedPostsResponse, err error) {
	post, err := m.Impl.PublishedPosts()
	if err != nil {
		return &protodef.SitePublishedPostsResponse{}, err
	}

	p, err := ArrayToProtobufStruct(post)
	return &protodef.SitePublishedPostsResponse{
		Posts: p,
	}, err
}

// PublishedPages handler.
func (m *GRPCSiteServer) PublishedPages(ctx context.Context, req *protodef.Empty) (
	resp *protodef.SitePublishedPagesResponse, err error) {
	post, err := m.Impl.PublishedPages()
	if err != nil {
		return &protodef.SitePublishedPagesResponse{}, err
	}

	p, err := ArrayToProtobufStruct(post)
	return &protodef.SitePublishedPagesResponse{
		Posts: p,
	}, err
}

// PostBySlug handler.
func (m *GRPCSiteServer) PostBySlug(ctx context.Context, req *protodef.SitePostBySlugRequest) (
	resp *protodef.SitePostBySlugResponse, err error) {
	post, err := m.Impl.PostBySlug(req.Slug)
	if err != nil {
		return &protodef.SitePostBySlugResponse{}, err
	}

	p, err := ObjectToProtobufStruct(post)
	return &protodef.SitePostBySlugResponse{
		Post: p,
	}, err
}

// PostByID handler.
func (m *GRPCSiteServer) PostByID(ctx context.Context, req *protodef.SitePostByIDRequest) (
	resp *protodef.SitePostByIDResponse, err error) {
	post, err := m.Impl.PostByID(req.Id)
	if err != nil {
		return &protodef.SitePostByIDResponse{}, err
	}

	p, err := ObjectToProtobufStruct(post)
	return &protodef.SitePostByIDResponse{
		Post: p,
	}, err
}

// DeletePostByID handler.
func (m *GRPCSiteServer) DeletePostByID(ctx context.Context, req *protodef.SiteDeletePostByIDRequest) (
	resp *protodef.Empty, err error) {
	err = m.Impl.DeletePostByID(req.Id)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, nil
}

// PluginNeighborRoutesList handler.
func (m *GRPCSiteServer) PluginNeighborRoutesList(ctx context.Context, req *protodef.SitePluginNeighborRoutesListRequest) (
	resp *protodef.SitePluginNeighborRoutesListResponse, err error) {
	routes, err := m.Impl.PluginNeighborRoutesList(req.Pluginname)
	if err != nil {
		return &protodef.SitePluginNeighborRoutesListResponse{}, err
	}

	r, err := ArrayToProtobufStruct(routes)
	return &protodef.SitePluginNeighborRoutesListResponse{
		Routes: r,
	}, err
}

// UserPersist handler.
func (m *GRPCSiteServer) UserPersist(ctx context.Context, req *protodef.SiteUserPersistRequest) (
	resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.UserPersist(c.Request, req.Persist)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, nil
}

// UserLogin handler.
func (m *GRPCSiteServer) UserLogin(ctx context.Context, req *protodef.SiteUserLoginRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.UserLogin(c.Request, req.Username)
	return &protodef.Empty{}, err
}

// AuthenticatedUser handler.
func (m *GRPCSiteServer) AuthenticatedUser(ctx context.Context, req *protodef.SiteAuthenticatedUserRequest) (resp *protodef.SiteAuthenticatedUserResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.SiteAuthenticatedUserResponse{}, err
	}

	username, err := m.Impl.AuthenticatedUser(c.Request)
	return &protodef.SiteAuthenticatedUserResponse{
		Username: username,
	}, err
}

// UserLogout handler.
func (m *GRPCSiteServer) UserLogout(ctx context.Context, req *protodef.SiteUserLogoutRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.UserLogout(c.Request)
	return &protodef.Empty{}, err
}

// LogoutAllUsers handler.
func (m *GRPCSiteServer) LogoutAllUsers(ctx context.Context, req *protodef.SiteLogoutAllUsersRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.LogoutAllUsers(c.Request)
	return &protodef.Empty{}, err
}

// SetCSRF handler.
func (m *GRPCSiteServer) SetCSRF(ctx context.Context, req *protodef.SiteSetCSRFRequest) (resp *protodef.SiteSetCSRFResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.SiteSetCSRFResponse{}, err
	}

	token := m.Impl.SetCSRF(c.Request)
	return &protodef.SiteSetCSRFResponse{
		Token: token,
	}, nil
}

// CSRF handler.
func (m *GRPCSiteServer) CSRF(ctx context.Context, req *protodef.SiteCSRFRequest) (resp *protodef.SiteCSRFResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.SiteCSRFResponse{
			Valid: false,
		}, err
	}

	valid := m.Impl.CSRF(c.Request, req.Token)
	return &protodef.SiteCSRFResponse{
		Valid: valid,
	}, nil
}

// SessionValue handler.
func (m *GRPCSiteServer) SessionValue(ctx context.Context, req *protodef.SiteSessionValueRequest) (resp *protodef.SiteSessionValueResponse, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.SiteSessionValueResponse{
			Value: "",
		}, err
	}

	val := m.Impl.SessionValue(c.Request, req.Name)
	return &protodef.SiteSessionValueResponse{
		Value: val,
	}, nil
}

// SetSessionValue handler.
func (m *GRPCSiteServer) SetSessionValue(ctx context.Context, req *protodef.SiteSetSessionValueRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	err = m.Impl.SetSessionValue(c.Request, req.Name, req.Value)
	return &protodef.Empty{}, err
}

// DeleteSessionValue handler.
func (m *GRPCSiteServer) DeleteSessionValue(ctx context.Context, req *protodef.SiteDeleteSessionValueRequest) (resp *protodef.Empty, err error) {
	c := m.reqmap.Load(req.Requestid)
	if c == nil {
		return &protodef.Empty{}, err
	}

	m.Impl.DeleteSessionValue(c.Request, req.Name)
	return &protodef.Empty{}, nil
}

// PluginNeighborSettingsList handler.
func (m *GRPCSiteServer) PluginNeighborSettingsList(ctx context.Context, req *protodef.SitePluginNeighborSettingsListRequest) (resp *protodef.SitePluginNeighborSettingsListResponse, err error) {
	settings, err := m.Impl.PluginNeighborSettingsList(req.Pluginname)
	if err != nil {
		return &protodef.SitePluginNeighborSettingsListResponse{}, err
	}

	arr, err := ArrayToProtobufStruct(settings)
	return &protodef.SitePluginNeighborSettingsListResponse{
		Settings: arr,
	}, err
}

// SetPluginSetting handler.
func (m *GRPCSiteServer) SetPluginSetting(ctx context.Context, req *protodef.SiteSetPluginSettingRequest) (resp *protodef.Empty, err error) {
	err = m.Impl.SetPluginSetting(req.Settingname, req.Value)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, nil
}

// PluginSettingBool handler.
func (m *GRPCSiteServer) PluginSettingBool(ctx context.Context, req *protodef.SitePluginSettingBoolRequest) (resp *protodef.SitePluginSettingBoolResponse, err error) {
	val, err := m.Impl.PluginSettingBool(req.Fieldname)
	if err != nil {
		return &protodef.SitePluginSettingBoolResponse{
			Value: false,
		}, err
	}

	return &protodef.SitePluginSettingBoolResponse{
		Value: val,
	}, nil
}

// PluginSettingString handler.
func (m *GRPCSiteServer) PluginSettingString(ctx context.Context, req *protodef.SitePluginSettingStringRequest) (resp *protodef.SitePluginSettingStringResponse, err error) {
	val, err := m.Impl.PluginSettingString(req.Fieldname)
	if err != nil {
		return &protodef.SitePluginSettingStringResponse{
			Value: "",
		}, err
	}

	return &protodef.SitePluginSettingStringResponse{
		Value: val,
	}, nil
}

// PluginSetting handler.
func (m *GRPCSiteServer) PluginSetting(ctx context.Context, req *protodef.SitePluginSettingRequest) (resp *protodef.SitePluginSettingResponse, err error) {
	setting, err := m.Impl.PluginSetting(req.Fieldname)
	if err != nil {
		return &protodef.SitePluginSettingResponse{}, err
	}

	val, err := InterfaceToProtobufAny(setting)
	return &protodef.SitePluginSettingResponse{
		Value: val,
	}, err
}

// SetNeighborPluginSetting handler.
func (m *GRPCSiteServer) SetNeighborPluginSetting(ctx context.Context, req *protodef.SiteSetNeighborPluginSettingRequest) (resp *protodef.Empty, err error) {
	err = m.Impl.SetNeighborPluginSetting(req.Pluginname, req.Settingname, req.Settingvalue)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, nil
}

// NeighborPluginSettingString handler.
func (m *GRPCSiteServer) NeighborPluginSettingString(ctx context.Context, req *protodef.SiteNeighborPluginSettingStringRequest) (resp *protodef.SiteNeighborPluginSettingStringResponse, err error) {
	val, err := m.Impl.NeighborPluginSettingString(req.Pluginname, req.Fieldname)
	if err != nil {
		return &protodef.SiteNeighborPluginSettingStringResponse{
			Value: "",
		}, err
	}

	return &protodef.SiteNeighborPluginSettingStringResponse{
		Value: val,
	}, nil
}

// NeighborPluginSetting handler.
func (m *GRPCSiteServer) NeighborPluginSetting(ctx context.Context, req *protodef.SiteNeighborPluginSettingRequest) (resp *protodef.SiteNeighborPluginSettingResponse, err error) {
	setting, err := m.Impl.NeighborPluginSetting(req.Pluginname, req.Fieldname)
	if err != nil {
		return &protodef.SiteNeighborPluginSettingResponse{}, err
	}

	val, err := InterfaceToProtobufAny(setting)
	return &protodef.SiteNeighborPluginSettingResponse{
		Value: val,
	}, err
}

// PluginTrusted handler.
func (m *GRPCSiteServer) PluginTrusted(ctx context.Context, req *protodef.SitePluginTrustedRequest) (resp *protodef.SitePluginTrustedResponse, err error) {
	trusted, err := m.Impl.PluginTrusted(req.Pluginname)
	if err != nil {
		return &protodef.SitePluginTrustedResponse{
			Trusted: false,
		}, err
	}

	return &protodef.SitePluginTrustedResponse{
		Trusted: trusted,
	}, err
}

// SetTitle handler.
func (m *GRPCSiteServer) SetTitle(ctx context.Context, req *protodef.SiteSetTitleRequest) (resp *protodef.Empty, err error) {
	err = m.Impl.SetTitle(req.Title)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, err
}

// Title handler.
func (m *GRPCSiteServer) Title(ctx context.Context, req *protodef.Empty) (resp *protodef.SiteTitleResponse, err error) {
	title, err := m.Impl.Title()
	if err != nil {
		return &protodef.SiteTitleResponse{
			Title: "",
		}, err
	}

	return &protodef.SiteTitleResponse{
		Title: title,
	}, err
}

// SetScheme handler.
func (m *GRPCSiteServer) SetScheme(ctx context.Context, req *protodef.SiteSetSchemeRequest) (resp *protodef.Empty, err error) {
	err = m.Impl.SetScheme(req.Scheme)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, err
}

// Scheme handler.
func (m *GRPCSiteServer) Scheme(ctx context.Context, req *protodef.Empty) (resp *protodef.SiteSchemeResponse, err error) {
	scheme, err := m.Impl.Scheme()
	if err != nil {
		return &protodef.SiteSchemeResponse{
			Scheme: "",
		}, err
	}

	return &protodef.SiteSchemeResponse{
		Scheme: scheme,
	}, err
}

// SetURL handler.
func (m *GRPCSiteServer) SetURL(ctx context.Context, req *protodef.SiteSetURLRequest) (resp *protodef.Empty, err error) {
	err = m.Impl.SetURL(req.Url)
	if err != nil {
		return &protodef.Empty{}, err
	}

	return &protodef.Empty{}, err
}

// URL handler.
func (m *GRPCSiteServer) URL(ctx context.Context, req *protodef.Empty) (resp *protodef.SiteURLResponse, err error) {
	URL, err := m.Impl.URL()
	if err != nil {
		return &protodef.SiteURLResponse{
			Url: "",
		}, err
	}

	return &protodef.SiteURLResponse{
		Url: URL,
	}, err
}

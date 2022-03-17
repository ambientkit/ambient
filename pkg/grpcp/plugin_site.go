package grpcp

import (
	"context"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/secureconfig"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCSitePlugin is the plugin side implementation of secure site.
type GRPCSitePlugin struct {
	client protodef.SiteClient
	Log    ambient.Logger
}

// Error handler.
func (c *GRPCSitePlugin) Error(siteError error) error {
	return secureconfig.Error(siteError)
}

// Load handler.
func (c *GRPCSitePlugin) Load() error {
	_, err := c.client.Load(context.Background(), &protodef.Empty{})
	return ErrorHandler(err)
}

// Authorized handler.
func (c *GRPCSitePlugin) Authorized(grant ambient.Grant) bool {
	resp, err := c.client.Authorized(context.Background(), &protodef.SiteAuthorizedRequest{
		Grant: string(grant),
	})
	if err != nil {
		c.Log.Error("grpc-plugin: site.Authorized error: %v", err.Error())
	}
	return resp.Authorized
}

// NeighborPluginGrantList handler.
func (c *GRPCSitePlugin) NeighborPluginGrantList(pluginName string) ([]ambient.GrantRequest, error) {
	resp, err := c.client.NeighborPluginGrantList(context.Background(), &protodef.SiteNeighborPluginGrantListRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return []ambient.GrantRequest{}, ErrorHandler(err)
	}

	arr := make([]ambient.GrantRequest, 0)
	for _, v := range resp.Grants {
		arr = append(arr, ambient.GrantRequest{
			Grant:       ambient.Grant(v.Grant),
			Description: v.Description,
		})
	}

	return arr, nil
}

// NeighborPluginGrants handler.
func (c *GRPCSitePlugin) NeighborPluginGrants(pluginName string) (map[ambient.Grant]bool, error) {
	resp, err := c.client.NeighborPluginGrants(context.Background(), &protodef.SiteNeighborPluginGrantsRequest{
		Pluginname: pluginName,
	})
	if err != nil {
		return make(map[ambient.Grant]bool), ErrorHandler(err)
	}

	sm, err := ProtobufStructToGrantBoolMap(resp.Grants)
	if err != nil {
		return make(map[ambient.Grant]bool), ErrorHandler(err)
	}

	return sm, nil
}

// NeighborPluginGranted handler.
func (c *GRPCSitePlugin) NeighborPluginGranted(pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginGranted(context.Background(), &protodef.SiteNeighborPluginGrantedRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// NeighborPluginRequestedGrant handler.
func (c *GRPCSitePlugin) NeighborPluginRequestedGrant(pluginName string, grantName ambient.Grant) (bool, error) {
	resp, err := c.client.NeighborPluginRequestedGrant(context.Background(), &protodef.SiteNeighborPluginRequestedGrantRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
	})
	if err != nil {
		return false, ErrorHandler(err)
	}

	return resp.Granted, nil
}

// SetNeighborPluginGrant handler.
func (c *GRPCSitePlugin) SetNeighborPluginGrant(pluginName string, grantName ambient.Grant, granted bool) error {
	_, err := c.client.SetNeighborPluginGrant(context.Background(), &protodef.SiteSetNeighborPluginGrantRequest{
		Pluginname: pluginName,
		Grant:      string(grantName),
		Granted:    granted,
	})
	if err != nil {
		return ErrorHandler(err)
	}

	return nil
}

// Plugins handler.
func (c *GRPCSitePlugin) Plugins() (map[string]ambient.PluginData, error) {
	resp, err := c.client.Plugins(context.Background(), &protodef.Empty{})
	if err != nil {
		return make(map[string]ambient.PluginData), ErrorHandler(err)
	}

	sm, err := ProtobufStructToPluginDataMap(resp.Plugindata)
	if err != nil {
		return make(map[string]ambient.PluginData), ErrorHandler(err)
	}

	return sm, nil
}

/////////////////////////////////////////////////////

// UserLogin handler.
func (c *GRPCSitePlugin) UserLogin(r *http.Request, username string) error {
	_, err := c.client.UserLogin(context.Background(), &protodef.SiteUserLoginRequest{
		Username:  username,
		Requestid: requestID(r),
	})
	return err
}

// AuthenticatedUser handler.
func (c *GRPCSitePlugin) AuthenticatedUser(r *http.Request) (string, error) {
	out, err := c.client.AuthenticatedUser(context.Background(), &protodef.SiteAuthenticatedUserRequest{
		Requestid: requestID(r),
	})
	if err != nil {
		return "", ErrorHandler(err)
	}

	return out.Username, err
}

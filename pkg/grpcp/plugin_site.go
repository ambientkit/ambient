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
		return make(map[ambient.Grant]bool), err
	}

	sm, err := ProtobufStructToGrantBoolMap(resp.Grants)
	if err != nil {
		return make(map[ambient.Grant]bool), err
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

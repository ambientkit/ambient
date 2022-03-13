package grpcp

import (
	"context"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCSitePlugin is the plugin side implementation of secure site.
type GRPCSitePlugin struct {
	client protodef.SiteClient
	Log    ambient.Logger
}

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
	return out.Username, err
}

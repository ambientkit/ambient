package grpcp

import (
	"context"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCSitePlugin .
type GRPCSitePlugin struct {
	client protodef.SiteClient
	Log    ambient.Logger
}

// requestID returns the request ID from the request context.
func (c *GRPCSitePlugin) requestID(r *http.Request) string {
	val := r.Context().Value(ambientRequestID)
	if val == nil {
		val = ""
	}
	return val.(string)
}

// UserLogin handler.
func (c *GRPCSitePlugin) UserLogin(r *http.Request, username string) error {
	_, err := c.client.UserLogin(context.Background(), &protodef.SiteUserLoginRequest{
		Username:  username,
		Requestid: c.requestID(r),
	})
	return err
}

// AuthenticatedUser handler.
func (c *GRPCSitePlugin) AuthenticatedUser(r *http.Request) (string, error) {
	out, err := c.client.AuthenticatedUser(context.Background(), &protodef.SiteAuthenticatedUserRequest{
		Requestid: c.requestID(r),
	})
	return out.Username, err
}

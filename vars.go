package ambient

import (
	"errors"
)

var (
	// ErrPluginNotFound returns a plugin that is not found.
	ErrPluginNotFound = errors.New("plugin name not found")

	// DisallowedPluginNames is a list of disallowed plugin names.
	DisallowedPluginNames = map[string]bool{
		"plugin":  false,
		"plugins": false,
		"ambient": false,
		"amb":     false,
	}
)

// GrantRequest represents a plugin grant request.
type GrantRequest struct {
	Grant       Grant
	Description string
}

// Route is a route for a router.
type Route struct {
	Method string
	Path   string
}

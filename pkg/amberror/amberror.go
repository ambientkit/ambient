// Package amberror has shared errors across the application.
package amberror

import "errors"

var (
	// ErrAccessDenied is when access is not allowed to the data item.
	ErrAccessDenied = errors.New("access denied to the data item")
	// ErrNotFound is when an item is not found.
	ErrNotFound = errors.New("item was not found")
	// ErrPluginNotFound returns a plugin that is not found.
	ErrPluginNotFound = errors.New("plugin name not found")
	// ErrGrantNotRequested is when a grant is attempted to be enabled on a
	// plugin, but the plugin didn't explicitly request the grant so don't allow
	// it.
	ErrGrantNotRequested = errors.New("request does not exist for the grant")
	// ErrSettingNotSpecified is when a setting is attempted to be set on a
	// plugin, but the plugin didn't explicity specify it as a setting.
	ErrSettingNotSpecified = errors.New("setting does not exist for the plugin")
)

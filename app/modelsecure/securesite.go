package modelsecure

import (
	"errors"
	"log"

	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/model"
)

var (
	// ErrAccessDenied is when access is not allowed to the data item.
	ErrAccessDenied = errors.New("access denied to the data item")
)

// SecureSite is a secure data access for the site.
type SecureSite struct {
	pluginName string
	storage    *datastorage.Storage
	grants     map[string]bool
}

// NewSecureSite -
func NewSecureSite(pluginName string, storage *datastorage.Storage, grants map[string]bool) SecureSite {
	return SecureSite{
		pluginName: pluginName,
		storage:    storage,
		grants:     grants,
	}
}

// Authorized determines if the current context has access.
func (ss SecureSite) Authorized(grant string) bool {
	//return true
	if allowed, ok := ss.grants[grant]; ok && allowed {
		return true
	}

	log.Printf("denied plugin (%v) access to the data action: %v\n", ss.pluginName, grant)

	return false
}

// Title returns the title.
func (ss SecureSite) Title() (string, error) {
	grant := "site.title:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Title, nil
}

// SetTitle sets the title.
func (ss SecureSite) SetTitle(title string) error {
	grant := "site.title:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Title = title

	return ss.storage.Save()
}

// Plugins returns the plugin list.
func (ss SecureSite) Plugins() (map[string]model.PluginSettings, error) {
	grant := "site.plugins:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.Plugins, nil
}

// Plugins returns the plugin list.
func (ss SecureSite) DeletePlugin(name string) error {
	grant := "site.plugins:deleteone"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.Plugins, name)

	return ss.storage.Save()
}

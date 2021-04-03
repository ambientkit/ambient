// Package scssession provides session capability with scs
// for an Ambient application.
package scssession

import (
	"embed"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/ambient/app/core"
	"github.com/josephspurrier/ambient/app/lib/websession"
)

//go:embed *
var assets embed.FS

// Plugin represents an Ambient plugin.
type Plugin struct {
	*core.PluginMeta
	*core.Toolkit

	sessionManager *scs.SessionManager
	sess           *websession.Session
}

// New returns a new scssession plugin.
func New() *Plugin {
	return &Plugin{
		PluginMeta: &core.PluginMeta{
			Name:       "scssession",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *core.Toolkit) error {
	p.Toolkit = toolkit

	return nil
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.sessionManager.LoadAndSave,
	}
}

// SessionManager return the session manager.
func (p *Plugin) SessionManager(ss core.SessionStorer, secretKey string) (core.ISession, error) {
	// Set up the session storage provider.
	en := websession.NewEncryptedStorage(secretKey)
	store, err := websession.NewJSONSession(ss, en)
	if err != nil {
		return nil, err
	}

	sessionName := "session"

	p.sessionManager = scs.New()
	p.sessionManager.Lifetime = 24 * time.Hour
	p.sessionManager.Cookie.Persist = false
	p.sessionManager.Store = store
	p.sess = websession.New(sessionName, p.sessionManager)

	return p.sess, nil
}

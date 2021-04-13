// Package scssession provides session capability with scs
// for an Ambient application.
package scssession

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/scssession/websession"
)

// Plugin represents an Ambient plugin.
type Plugin struct {
	*ambient.PluginBase
	*ambient.Toolkit

	sessionManager *scs.SessionManager
	sess           *websession.Session
}

// New returns a new scssession plugin.
func New() *Plugin {
	return &Plugin{
		PluginBase: &ambient.PluginBase{},
	}
}

// PluginName returns the plugin name.
func (p *Plugin) PluginName() string {
	return "scssession"
}

// PluginVersion returns the plugin version.
func (p *Plugin) PluginVersion() string {
	return "1.0.0"
}

// Enable accepts the toolkit.
func (p *Plugin) Enable(toolkit *ambient.Toolkit) error {
	p.Toolkit = toolkit

	return nil
}

// Middleware returns router middleware.
func (p *Plugin) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		p.sessionManager.LoadAndSave,
	}
}

// SessionManager returns the session manager.
func (p *Plugin) SessionManager(logger ambient.ILogger, ss ambient.SessionStorer) (ambient.IAppSession, error) {
	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("environment variable missing: %v", "AMB_SESSION_KEY")
	}

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

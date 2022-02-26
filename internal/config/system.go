package config

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ambientkit/ambient"
)

//go:generate go run github.com/vburenin/ifacemaker -f *.go -s PluginSystem -i PluginSystem -p ambient -o ../../gen_pluginsystem.go -y "PluginSystem provides config functions." -c "Code generated by ifacemaker. DO NOT EDIT."

var (
	// ErrAccessDenied is when access is not allowed to the data item.
	ErrAccessDenied = errors.New("access denied to the data item")
	// ErrNotFound is when an item is not found.
	ErrNotFound = errors.New("item was not found")
	// ErrGrantNotRequested is when a grant is attempted to be enabled on a
	// plugin, but the plugin didn't explicitly request the grant so don't allow
	// it.
	ErrGrantNotRequested = errors.New("request does not exist for the grant")
	// ErrSettingNotSpecified is when a setting is attempted to be set on a
	// plugin, but the plugin didn't explicity specify it as a setting.
	ErrSettingNotSpecified = errors.New("setting does not exist for the plugin")
)

// PluginSystem represents loaded plugins.
type PluginSystem struct {
	log     ambient.AppLogger
	storage *Storage

	names           []string
	middlewareNames []string
	router          ambient.RouterPlugin
	templateEngine  ambient.TemplateEnginePlugin
	sessionManager  ambient.SessionManagerPlugin
	plugins         map[string]ambient.Plugin
	trusted         map[string]bool

	routes map[string][]ambient.Route
}

// NewPluginSystem returns a plugin system.
func NewPluginSystem(log ambient.AppLogger, storage *Storage, arr *ambient.PluginLoader) (*PluginSystem, error) {
	// Get a list of plugin names to maintain order.
	names := make([]string, 0)
	middlewareNames := make([]string, 0)
	plugins := make(map[string]ambient.Plugin)
	shouldSave := false

	// Load the middleware.
	for _, p := range arr.Middleware {
		// Don't allow certain plugin names.
		if allowed, ok := ambient.DisallowedPluginNames[p.PluginName()]; ok && !allowed {
			return nil, fmt.Errorf("ambient: plugin name not allowed: %v", p.PluginName())
		}

		save, err := loadPlugin(log, p, plugins, storage)
		if err != nil {
			return nil, err
		} else if save {
			shouldSave = true
		}

		names = append(names, p.PluginName())
		middlewareNames = append(middlewareNames, p.PluginName())
	}

	// Load the plugins.
	for _, p := range arr.Plugins {
		// Don't allow certain plugin names.
		if allowed, ok := ambient.DisallowedPluginNames[p.PluginName()]; ok && !allowed {
			return nil, fmt.Errorf("ambient: plugin name not allowed: %v", p.PluginName())
		}

		save, err := loadPlugin(log, p, plugins, storage)
		if err != nil {
			return nil, err
		} else if save {
			shouldSave = true
		}

		names = append(names, p.PluginName())
	}

	if shouldSave {
		err := storage.Save()
		if err != nil {
			return nil, err
		}
	}

	return &PluginSystem{
		log:     log,
		storage: storage,

		names:           names,
		middlewareNames: middlewareNames,
		router:          arr.Router,
		templateEngine:  arr.TemplateEngine,
		sessionManager:  arr.SessionManager,
		trusted:         arr.TrustedPlugins,

		plugins: plugins,
		routes:  make(map[string][]ambient.Route),
	}, nil
}

// SessionManager returns the session manager.
func (p *PluginSystem) SessionManager() ambient.SessionManagerPlugin {
	return p.sessionManager
}

// TemplateEngine returns the template engine.
func (p *PluginSystem) TemplateEngine() ambient.TemplateEnginePlugin {
	return p.templateEngine
}

// Router returns the router.
func (p *PluginSystem) Router() ambient.RouterPlugin {
	return p.router
}

// StorageManager returns the storage manager.
func (p *PluginSystem) StorageManager() ambient.Storage {
	return p.storage
}

func loadPlugin(log ambient.AppLogger, p ambient.Plugin, plugins map[string]ambient.Plugin, storage *Storage) (shouldSave bool, err error) {
	// TODO: Need to make sure the name matches a certain format. All lowercase. No symbols.

	// Ensure a plugin can't be loaded twice or two plugins with the same
	// names can't both be loaded.
	if _, found := plugins[p.PluginName()]; found {
		return false, fmt.Errorf("found a duplicate plugin: %v", p.PluginName())
	}

	plugins[p.PluginName()] = p

	// Determine if plugin if found in app config.
	pluginData, ok := storage.site.PluginStorage[p.PluginName()]
	if !ok {
		storage.site.PluginStorage[p.PluginName()] = newPluginData(p.PluginVersion())
		return true, nil
	}

	// Detect plugin version change.
	if pluginData.Version != p.PluginVersion() {
		log.Info("ambient: detected plugin (%v) version change from (%v) to: %v", p.PluginName(), pluginData.Version, p.PluginVersion())
		pluginData.Version = p.PluginVersion()
		storage.site.PluginStorage[p.PluginName()] = pluginData
		return true, nil
	}

	return false, nil
}

// newPluginData returns new PluginData.
func newPluginData(version string) ambient.PluginData {
	return ambient.PluginData{
		Enabled:  false,
		Version:  version,
		Grants:   make(ambient.PluginGrants),
		Settings: make(ambient.PluginSettings),
	}
}

// Load will load the storage and return an error if one occurs.
func (p *PluginSystem) Load() error {
	return p.storage.Load()
}

// Save will save the storage and return an error if one occurs.
func (p *PluginSystem) Save() error {
	return p.storage.Save()
}

// InitializePlugin will initialize the plugin in the storage and will return
// an error if one occurs.
func (p *PluginSystem) InitializePlugin(pluginName string, pluginVersion string) error {
	_, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.storage.site.PluginStorage[pluginName] = newPluginData(pluginVersion)
		return p.storage.Save()
	}

	return nil
}

// RemovePlugin will delete the plugin from the storage and will return
// an error if one occurs.
func (p *PluginSystem) RemovePlugin(pluginName string) error {
	_, ok := p.storage.site.PluginStorage[pluginName]
	if ok {
		delete(p.storage.site.PluginStorage, pluginName)
		return p.storage.Save()
	}

	return nil
}

// Names returns a list of plugin names.
func (p *PluginSystem) Names() []string {
	return p.names
}

// MiddlewareNames returns a list of middleware plugin names.
func (p *PluginSystem) MiddlewareNames() []string {
	return p.middlewareNames
}

// TrustedPluginNames returns a list of sorted trusted names.
func (p *PluginSystem) TrustedPluginNames() []string {
	names := make([]string, 0)
	for name, trust := range p.trusted {
		if trust {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// Trusted returns if a plugin is trusted.
func (p *PluginSystem) Trusted(pluginName string) bool {
	trust, found := p.trusted[pluginName]
	if !found {
		return false
	}

	return trust
}

// Routes returns a list of plugin routes.
func (p *PluginSystem) Routes(pluginName string) []ambient.Route {
	routes, found := p.routes[pluginName]
	if !found {
		return make([]ambient.Route, 0)
	}

	return routes
}

// Plugins returns the plugin list.
func (p *PluginSystem) Plugins() map[string]ambient.PluginData {
	return p.storage.site.PluginStorage
}

// Plugin returns a plugin by name.
func (p *PluginSystem) Plugin(name string) (ambient.Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, ambient.ErrPluginNotFound
	}

	return plugin, nil
}

// PluginData returns a plugin data by name.
func (p *PluginSystem) PluginData(name string) (ambient.PluginData, error) {
	plugin, ok := p.storage.site.PluginStorage[name]
	if !ok {
		return ambient.PluginData{}, ambient.ErrPluginNotFound
	}

	return plugin, nil
}

// Enabled returns if the plugin is enabled or not. If it cannot be found, it
// will still return false.
func (p *PluginSystem) Enabled(name string) bool {
	data, ok := p.storage.site.PluginStorage[name]
	if !ok {
		p.log.Debug("pluginsystem.enabled: could not find plugin: %v", name)
		return false
	}

	return data.Enabled
}

// SetEnabled sets a plugin as enabled or not.
func (p *PluginSystem) SetEnabled(pluginName string, enabled bool) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setenabled: could not find plugin: %v", pluginName)
		return ErrNotFound
	}

	data.Enabled = enabled
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// GrantRequests returns a list of grant requests.
func (p *PluginSystem) GrantRequests(pluginName string, grant ambient.Grant) ([]ambient.GrantRequest, error) {
	plugin, err := p.Plugin(pluginName)
	if err != nil {
		return nil, err
	}

	return plugin.GrantRequests(), nil
}

// Authorized returns whether a plugin is inherited granted for a plugin.
func (p *PluginSystem) Authorized(pluginName string, grant ambient.Grant) bool {
	// Always allow ambient system to get full access.
	if pluginName == "ambient" {
		p.log.Debug("pluginsystem: granted system (%v) GrantAll access to the data item for grant: %v", "ambient", grant)
		return true
	}

	// If the grant was found, then allow access.
	if granted := p.Granted(pluginName, grant); granted {
		p.log.Debug("pluginsystem: granted plugin (%v) access to the data item for grant: %v", pluginName, grant)
		return true
	}

	p.log.Warn("pluginsystem: denied plugin (%v) access to the data item, requires grant: %v", pluginName, grant)

	return false
}

// Granted returns whether a plugin is explicitly granted for a plugin.
func (p *PluginSystem) Granted(pluginName string, grant ambient.Grant) bool {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.granted: could not find plugin: %v", pluginName)
		return false
	}

	granted, found := data.Grants[grant]
	if !found {
		p.log.Debug("pluginsystem.granted: could not find grant for plugin (%v): %v", pluginName, grant)
		return false
	}

	return granted
}

// SetGrant sets a plugin grant.
func (p *PluginSystem) SetGrant(pluginName string, grant ambient.Grant) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setgrant: could not find plugin: %v", pluginName)
		return ErrNotFound
	}

	data.Grants[grant] = true
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// RemoveGrant removes a plugin grant.
func (p *PluginSystem) RemoveGrant(pluginName string, grant ambient.Grant) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.removegrant: could not find plugin: %v", pluginName)
		return ErrNotFound
	}

	delete(data.Grants, grant)
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// SetSetting sets a plugin setting.
func (p *PluginSystem) SetSetting(pluginName string, settingName string, value interface{}) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setsettings: could not find plugin: %v", pluginName)
		return ErrNotFound
	}

	data.Settings[settingName] = value
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// Setting returns a setting value.
func (p *PluginSystem) Setting(pluginName string, settingName string) (interface{}, error) {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setting: could not find plugin: %v", pluginName)
		return nil, ErrNotFound
	}

	value, ok := data.Settings[settingName]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

// SettingDefault returns a setting default for a setting.
func (p *PluginSystem) SettingDefault(pluginName string, settingName string) (interface{}, error) {
	plugin, err := p.Plugin(pluginName)
	if err != nil {
		return nil, err
	}

	// TODO: this needs to be more efficient.
	fields := plugin.Settings()
	for _, field := range fields {
		if field.Name == settingName {
			return field.Default, nil
		}
	}

	return nil, nil
}

// SetRoute saves a route.
func (p *PluginSystem) SetRoute(pluginName string, route []ambient.Route) {
	p.routes[pluginName] = route
}
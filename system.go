package ambient

import (
	"errors"
	"fmt"
	"sort"
)

var (
	// ErrPluginNotFound returns a plugin that is not found.
	ErrPluginNotFound = errors.New("plugin name not found")

	disallowedPluginNames = map[string]bool{
		"plugin":  false,
		"plugins": false,
		"ambient": false,
		"amb":     false,
	}
)

// PluginSystem represents loaded plugins.
type PluginSystem struct {
	log     AppLogger
	storage *Storage

	names           []string
	middlewareNames []string
	router          RouterPlugin
	templateEngine  TemplateEnginePlugin
	sessionManager  SessionManagerPlugin
	plugins         map[string]Plugin
	trusted         map[string]bool

	routes map[string][]Route
}

// NewPluginSystem returns a plugin system.
func NewPluginSystem(log AppLogger, storage *Storage, arr *PluginLoader) (*PluginSystem, error) {
	// Get a list of plugin names to maintain order.
	names := make([]string, 0)
	middlewareNames := make([]string, 0)
	plugins := make(map[string]Plugin)
	shouldSave := false

	// Load the middleware.
	for _, p := range arr.Middleware {
		// Don't allow certain plugin names.
		if allowed, ok := disallowedPluginNames[p.PluginName()]; ok && !allowed {
			return nil, fmt.Errorf("ambient: plugin name not allowed: %v", p.PluginName())
		}

		save, err := loadPlugin(p, plugins, storage)
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
		if allowed, ok := disallowedPluginNames[p.PluginName()]; ok && !allowed {
			return nil, fmt.Errorf("ambient: plugin name not allowed: %v", p.PluginName())
		}

		save, err := loadPlugin(p, plugins, storage)
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
		routes:  make(map[string][]Route),
	}, nil
}

func loadPlugin(p Plugin, plugins map[string]Plugin, storage *Storage) (shouldSave bool, err error) {
	// TODO: Need to make sure the name matches a certain format. All lowercase. No symbols.

	// Ensure a plugin can't be loaded twice or two plugins with the same
	// names can't both be loaded.
	if _, found := plugins[p.PluginName()]; found {
		return false, fmt.Errorf("found a duplicate plugin: %v", p.PluginName())
	}

	plugins[p.PluginName()] = p

	_, ok := storage.site.PluginStorage[p.PluginName()]
	if !ok {
		storage.site.PluginStorage[p.PluginName()] = newPluginData(p.PluginVersion())
		return true, nil
	}

	return false, nil
}

// newPluginData returns new PluginData.
func newPluginData(version string) PluginData {
	return PluginData{
		Enabled:  false,
		Version:  version,
		Grants:   make(PluginGrants),
		Settings: make(PluginSettings),
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

// Plugins returns the plugin list.
func (p *PluginSystem) Plugins() map[string]PluginData {
	return p.storage.site.PluginStorage
}

// Plugin returns a plugin by name.
func (p *PluginSystem) Plugin(name string) (Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, ErrPluginNotFound
	}

	return plugin, nil
}

// PluginData returns a plugin data by name.
func (p *PluginSystem) PluginData(name string) (PluginData, error) {
	plugin, ok := p.storage.site.PluginStorage[name]
	if !ok {
		return PluginData{}, ErrPluginNotFound
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
func (p *PluginSystem) GrantRequests(pluginName string, grant Grant) ([]GrantRequest, error) {
	plugin, err := p.Plugin(pluginName)
	if err != nil {
		return nil, err
	}

	return plugin.GrantRequests(), nil
}

// Authorized returns whether a plugin is inherited granted for a plugin.
func (p *PluginSystem) Authorized(pluginName string, grant Grant) bool {
	// Always allow ambient plugin to get full access.
	if pluginName == "ambient" {
		p.log.Debug("pluginsystem: granted plugin (%v) GrantAll access to the data item for grant: %v", "ambient", grant)
		return true
	}

	// If has star, then allow all access.
	if granted := p.Granted(pluginName, GrantAll); granted {
		p.log.Debug("pluginsystem: granted plugin (%v) GrantAll access to the data item for grant: %v", pluginName, grant)
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
func (p *PluginSystem) Granted(pluginName string, grant Grant) bool {
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
func (p *PluginSystem) SetGrant(pluginName string, grant Grant) error {
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
func (p *PluginSystem) RemoveGrant(pluginName string, grant Grant) error {
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
func (p *PluginSystem) SetRoute(pluginName string, route []Route) {
	p.routes[pluginName] = route
}

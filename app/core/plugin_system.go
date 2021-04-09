package core

import "errors"

var (
	// ErrPluginNotFound returns a plugin that is not found.
	ErrPluginNotFound = errors.New("plugin name not found")
)

// PluginSystem represents loaded plugins.
type PluginSystem struct {
	log     ILogger
	storage *Storage

	names   []string
	plugins map[string]IPlugin
	routes  map[string]IRouteList
}

// NewPluginSystem returns a plugin system.
func NewPluginSystem(log ILogger, arr []IPlugin, storage *Storage) *PluginSystem {
	// Get a list of plugin names to maintain order.
	names := make([]string, 0)
	plugins := make(map[string]IPlugin)

	for _, p := range arr {
		names = append(names, p.PluginName())
		plugins[p.PluginName()] = p
	}

	return &PluginSystem{
		log:     log,
		storage: storage,

		names:   names,
		plugins: plugins,
	}
}

// Names returns a list of plugin names.
func (p *PluginSystem) initialize() PluginData {
	return PluginData{
		Enabled:  false,
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

// Names returns a list of plugin names.
func (p *PluginSystem) Names() []string {
	return p.names
}

// Plugin returns a plugin by name.
func (p *PluginSystem) Plugin(name string) (IPlugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, ErrPluginNotFound
	}

	return plugin, nil
}

// Enabled returns if the plugin is enabled or not. If it cannot be found, it
// will still return false.
func (p *PluginSystem) Enabled(name string) bool {
	data, ok := p.storage.Site.PluginStorage[name]
	if !ok {
		p.log.Debug("pluginsystem.enabled: could not find plugin: %v", name)
		return false
	}

	return data.Enabled
}

// SetEnabled sets a plugin as enabled or not.
func (p *PluginSystem) SetEnabled(pluginName string, enabled bool) error {
	data, ok := p.storage.Site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setenabled: could not find plugin: %v", pluginName)
		data = p.initialize()
	}

	data.Enabled = enabled
	p.storage.Site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// GrantRequests returns a list of grant requests.
func (p *PluginSystem) GrantRequests(pluginName string, grant Grant) ([]Grant, error) {
	plugin, err := p.Plugin(pluginName)
	if err != nil {
		return nil, err
	}

	return plugin.Grants(), nil
}

// Granted returns whether a plugin is granted for a plugin.
func (p *PluginSystem) Granted(pluginName string, grant Grant) bool {
	data, ok := p.storage.Site.PluginStorage[pluginName]
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
	data, ok := p.storage.Site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setgrant: could not find plugin: %v", pluginName)
		data = p.initialize()
	}

	data.Grants[grant] = true
	p.storage.Site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// RemoveGrant removes a plugin grant.
func (p *PluginSystem) RemoveGrant(pluginName string, grant Grant) error {
	data, ok := p.storage.Site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.removegrant: could not find plugin: %v", pluginName)
		data = p.initialize()
	}

	delete(data.Grants, grant)
	p.storage.Site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// SetSetting sets a plugin setting.
func (p *PluginSystem) SetSetting(pluginName string, settingName string, value interface{}) error {
	data, ok := p.storage.Site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setsettings: could not find plugin: %v", pluginName)
		data = p.initialize()
	}

	data.Settings[settingName] = value
	p.storage.Site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// Setting returns a setting value.
func (p *PluginSystem) Setting(pluginName string, settingName string) (interface{}, error) {
	data, ok := p.storage.Site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("pluginsystem.setting: could not find plugin: %v", pluginName)
		data = p.initialize()
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
	fields := plugin.Fields()
	for _, field := range fields {
		if field.Name == settingName {
			return field.Default, nil
		}
	}

	return nil, nil
}

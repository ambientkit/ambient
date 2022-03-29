package config

import (
	"fmt"
	"sort"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

//go:generate go run github.com/vburenin/ifacemaker -f *.go -s PluginSystem -i PluginSystem -p ambient -o ../../gen_pluginsystem.go -y "PluginSystem provides config functions." -c "Code generated by ifacemaker. DO NOT EDIT."

// PluginSystem represents loaded plugins.
type PluginSystem struct {
	log     ambient.AppLogger
	storage *Storage
	loader  *ambient.PluginLoader

	// pluginNames contains the list of plugins and middleware. It does not
	// include the: logger, storage, template engine, or session manager
	// (unless session manager has middleware) because they are not
	// configurable (they are not passed Toolkit). They aren't configurable
	// because they have to load before the rest of the plugin system.
	pluginNames []string
	// middlewareNames contains the list of only middleware.
	middlewareNames []string
	// middlewareNamesMap contains thel ist of middleware names for quick lookup.
	middlewareNamesMap map[string]bool
	// plugins is a map of plugins and plugins for quick lookup.
	plugins map[string]ambient.Plugin
	// grpcPlugins tracks a list of gRPC plugins vs standard plugins.
	grpcPlugins map[string]bool
	// routes contains a map of routes for quick lookup. It's used to show which
	// routes will be added to the router before a user enables a plugin. It's
	// useful for the plugin manager so people don't enable plugins blindly.
	routes map[string][]ambient.Route
}

// NewPluginSystem returns a plugin system.
func NewPluginSystem(log ambient.AppLogger, storage *Storage, loader *ambient.PluginLoader) (*PluginSystem, error) {
	ps := &PluginSystem{
		log:     log,
		storage: storage,
		loader:  loader,

		pluginNames:        make([]string, 0),
		middlewareNames:    make([]string, 0),
		middlewareNamesMap: make(map[string]bool),
		plugins:            make(map[string]ambient.Plugin),
		grpcPlugins:        make(map[string]bool),
		routes:             make(map[string][]ambient.Route),
	}

	// shouldSave is for efficiency so there is not saving on every plugin.
	shouldSave := false

	// Load the middleware.
	for _, p := range loader.Middleware {
		// Ensure a plugin can't be loaded twice or two plugins with the same
		// names can't both be loaded.
		if _, found := ps.plugins[p.PluginName()]; found {
			return nil, fmt.Errorf("found a duplicate plugin: %v", p.PluginName())
		}

		// Skip gRPC plugins.
		if p.PluginVersion() == "gRPC" {
			continue
		}

		// Else, load the standard plugin.
		save, err := ps.loadPlugin(p, true, false)
		if err != nil {
			return nil, err
		} else if save {
			shouldSave = true
		}
	}

	// Load the plugins.
	for _, p := range loader.Plugins {
		// Ensure a plugin can't be loaded twice or two plugins with the same
		// names can't both be loaded.
		if _, found := ps.plugins[p.PluginName()]; found {
			return nil, fmt.Errorf("found a duplicate plugin: %v", p.PluginName())
		}

		// Skip gRPC plugins.
		if p.PluginVersion() == "gRPC" {
			continue
		}

		// Else, load the standard plugin.
		save, err := ps.loadPlugin(p, false, false)
		if err != nil {
			return nil, err
		} else if save {
			shouldSave = true
		}
	}

	if shouldSave {
		err := storage.Save()
		if err != nil {
			return nil, err
		}
	}

	return ps, nil
}

// LoaderPlugins returns the loader plugins, these include initial gRPC plugins
// as well.
func (p *PluginSystem) LoaderPlugins() []ambient.Plugin {
	return p.loader.Plugins
}

// LoaderMiddleware returns the loader middleware, these include initial gRPC plugins
// as well.
func (p *PluginSystem) LoaderMiddleware() []ambient.MiddlewarePlugin {
	return p.loader.Middleware
}

// Plugins returns the map of plugins. This returns a map that is passed by
// reference and all the values are pointers so any changes to it will
// be reflected in the plugin system.
func (p *PluginSystem) Plugins() map[string]ambient.Plugin {
	return p.plugins
}

// SessionManager returns the session manager.
func (p *PluginSystem) SessionManager() ambient.SessionManagerPlugin {
	return p.loader.SessionManager
}

// TemplateEngine returns the template engine.
func (p *PluginSystem) TemplateEngine() ambient.TemplateEnginePlugin {
	return p.loader.TemplateEngine
}

// Router returns the router.
func (p *PluginSystem) Router() ambient.RouterPlugin {
	return p.loader.Router
}

// StorageManager returns the storage manager.
func (p *PluginSystem) StorageManager() ambient.Storage {
	return p.storage
}

// LoadPlugin loads a single plugin into the plugin system and saves the config.
func (p *PluginSystem) LoadPlugin(plugin ambient.Plugin, middleware bool, grpcPlugin bool) (err error) {
	shouldSave, err := p.loadPlugin(plugin, middleware, grpcPlugin)
	if err != nil {
		return err
	}

	// TODO: Add these so they are in the loader. These may be other work that
	// needs to happen as well. They have to be in the loader for the
	// grpcsystem to be able to revive them.
	// if middleware {
	// 	//p.loader.Middleware
	// } else {
	// 	// p.loader.Plugins
	// }

	if shouldSave {
		err = p.storage.Save()
	}

	return err
}

// loadPlugin adds a plugin to the plugin system and returns whether the config
// should be saved.
func (p *PluginSystem) loadPlugin(plugin ambient.Plugin, middleware bool, grpcPlugin bool) (shouldSave bool, err error) {
	// Validate plugin name and version.
	err = ambient.Validate(plugin)
	if err != nil {
		return false, err
	}
	name := plugin.PluginName()
	version := plugin.PluginVersion()

	isGRPC, found := p.grpcPlugins[plugin.PluginName()]
	if found {
		if grpcPlugin != isGRPC {
			return false, fmt.Errorf("cannot load the same plugin of two different types (gRPC/non-gRPC): %v", plugin.PluginName())
		}
	}

	// Determine if an old plugin is found already loaded.
	oldPlugin, exists := p.plugins[name]
	if exists {
		oldPlugin.Disable()
	}

	// Store the plugin.
	p.plugins[name] = plugin
	p.grpcPlugins[plugin.PluginName()] = grpcPlugin
	if !exists {
		p.pluginNames = append(p.pluginNames, plugin.PluginName())
	}
	if middleware {
		p.middlewareNamesMap[plugin.PluginName()] = true
		if !exists {
			p.middlewareNames = append(p.middlewareNames, plugin.PluginName())
		}
	}

	// Determine if plugin if found in app config.
	pluginData, ok := p.storage.site.PluginStorage[name]
	if !ok {
		p.storage.site.PluginStorage[name] = newPluginData(version)
		return true, nil
	}

	// Detect plugin version change.
	if pluginData.Version != version {
		p.log.Info("detected plugin (%v) version change from (%v) to: %v", name, pluginData.Version, version)
		pluginData.Version = version
		p.storage.site.PluginStorage[name] = pluginData
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
	// Make a copy to prevent order changing via sorting.
	out := make([]string, len(p.pluginNames))
	copy(out, p.pluginNames)
	return out
}

// MiddlewareNames returns a list of middleware plugin names.
func (p *PluginSystem) MiddlewareNames() []string {
	// Make a copy to prevent order changing via sorting.
	out := make([]string, len(p.middlewareNames))
	copy(out, p.middlewareNames)
	return out
}

// IsMiddleware returns if the plugin is middleware.
func (p *PluginSystem) IsMiddleware(name string) bool {
	if b, ok := p.middlewareNamesMap[name]; ok {
		return b
	}

	return false
}

// TrustedPluginNames returns a list of sorted trusted names.
func (p *PluginSystem) TrustedPluginNames() []string {
	names := make([]string, 0)
	for name, trust := range p.loader.TrustedPlugins {
		if trust {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// Trusted returns if a plugin is trusted.
func (p *PluginSystem) Trusted(pluginName string) bool {
	trust, found := p.loader.TrustedPlugins[pluginName]
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

// PluginsData returns the plugin data map.
func (p *PluginSystem) PluginsData() map[string]ambient.PluginData {
	// Create a new map so it doesn't copy by reference.
	m := make(map[string]ambient.PluginData)
	for k, v := range p.storage.site.PluginStorage {
		m[k] = v
	}
	return m
}

// Plugin returns a plugin by name.
func (p *PluginSystem) Plugin(name string) (ambient.Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, amberror.ErrPluginNotFound
	}

	return plugin, nil
}

// PluginData returns a plugin data by name.
func (p *PluginSystem) PluginData(name string) (ambient.PluginData, error) {
	plugin, ok := p.storage.site.PluginStorage[name]
	if !ok {
		return ambient.PluginData{}, amberror.ErrPluginNotFound
	}

	return plugin, nil
}

// Enabled returns if the plugin is enabled or not. If it cannot be found, it
// will still return false.
func (p *PluginSystem) Enabled(name string) bool {
	data, ok := p.storage.site.PluginStorage[name]
	if !ok {
		p.log.Debug("could not find plugin: %v", name)
		return false
	}

	return data.Enabled
}

// SetEnabled sets a plugin as enabled or not.
func (p *PluginSystem) SetEnabled(pluginName string, enabled bool) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return amberror.ErrNotFound
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
		p.log.Debug("granted system (%v) GrantAll access to the data item for grant: %v", "ambient", grant)
		return true
	}

	// If the grant was found, then allow access.
	if granted := p.Granted(pluginName, grant); granted {
		p.log.Debug("granted plugin (%v) access to the data item for grant: %v", pluginName, grant)
		return true
	}

	p.log.Warn("denied plugin (%v) access to the data item, requires grant: %v", pluginName, grant)

	return false
}

// Granted returns whether a plugin is explicitly granted for a plugin.
func (p *PluginSystem) Granted(pluginName string, grant ambient.Grant) bool {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return false
	}

	granted, found := data.Grants[grant]
	if !found {
		p.log.Debug("could not find grant for plugin (%v): %v", pluginName, grant)
		return false
	}

	return granted
}

// SetGrant sets a plugin grant.
func (p *PluginSystem) SetGrant(pluginName string, grant ambient.Grant) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return amberror.ErrNotFound
	}

	data.Grants[grant] = true
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// RemoveGrant removes a plugin grant.
func (p *PluginSystem) RemoveGrant(pluginName string, grant ambient.Grant) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return amberror.ErrNotFound
	}

	delete(data.Grants, grant)
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// SetSetting sets a plugin setting.
func (p *PluginSystem) SetSetting(pluginName string, settingName string, value interface{}) error {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return amberror.ErrNotFound
	}

	data.Settings[settingName] = value
	p.storage.site.PluginStorage[pluginName] = data

	return p.storage.Save()
}

// Setting returns a setting value.
func (p *PluginSystem) Setting(pluginName string, settingName string) (interface{}, error) {
	data, ok := p.storage.site.PluginStorage[pluginName]
	if !ok {
		p.log.Debug("could not find plugin: %v", pluginName)
		return nil, amberror.ErrNotFound
	}

	value, ok := data.Settings[settingName]
	if !ok {
		return nil, amberror.ErrNotFound
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

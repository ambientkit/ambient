package core

import "fmt"

// Plugins returns the plugin list.
func (ss *SecureSite) Plugins() (map[string]PluginSettings, error) {
	if !ss.Authorized(GrantSitePluginRead) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PluginSettings, nil
}

// DeletePlugin deletes a plugin.
func (ss *SecureSite) DeletePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginDelete) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.PluginSettings, name)

	return ss.storage.Save()
}

// EnablePlugin enables a plugin.
func (ss *SecureSite) EnablePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginEnable) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginSettings[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = true
	ss.storage.Site.PluginSettings[name] = plugin

	return ss.storage.Save()
}

// DisablePlugin disables a plugin.
func (ss *SecureSite) DisablePlugin(name string) error {
	if !ss.Authorized(GrantSitePluginDisable) {
		return ErrAccessDenied
	}

	plugin, ok := ss.storage.Site.PluginSettings[name]
	if !ok {
		return ErrNotFound
	}

	plugin.Enabled = false
	ss.storage.Site.PluginSettings[name] = plugin

	return ss.storage.Save()
}

// SetPluginField sets a variable for the plugin.
func (ss *SecureSite) SetPluginField(name string, value string) error {
	if !ss.Authorized(GrantPluginFieldWrite) {
		return ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[ss.pluginName]
	if !ok {
		fields = PluginFields{
			Fields: make(map[string]interface{}),
		}
	}

	fields.Fields[name] = value
	ss.storage.Site.PluginFields[ss.pluginName] = fields

	return ss.storage.Save()
}

func (ss *SecureSite) pluginField(pluginName string, fieldName string) (interface{}, error) {
	// See if the value is set.
	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if ok {
		value, ok := fields.Fields[fieldName]
		if ok {
			return value, nil
		}
	}

	// See if there is a default value.
	plugin, ok := ss.storage.Site.PluginSettings[pluginName]
	if ok {
		field, ok := plugin.Fields[fieldName]
		if ok {
			return field.Default, nil
		}
	}

	return "", nil
}

// PluginFieldBool returns a plugin field as a bool.
func (ss *SecureSite) PluginFieldBool(name string) (bool, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return false, ErrAccessDenied
	}

	value, err := ss.pluginField(ss.pluginName, name)

	return value == "true", err
}

// PluginFieldString gets a variable for the plugin as a string.
func (ss *SecureSite) PluginFieldString(fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.pluginField(ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// PluginField gets a variable for the plugin as an interface{}.
func (ss *SecureSite) PluginField(fieldName string) (interface{}, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.pluginField(ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// SetNeighborPluginField sets a variable for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginField(pluginName string, fieldName string, value string) error {
	if !ss.Authorized(GrantPluginNeighborfieldWrite) {
		return ErrAccessDenied
	}

	fields, ok := ss.storage.Site.PluginFields[pluginName]
	if !ok {
		fields = PluginFields{
			Fields: make(map[string]interface{}),
		}
	}

	fields.Fields[fieldName] = value
	ss.storage.Site.PluginFields[pluginName] = fields

	return ss.storage.Save()
}

// NeighborPluginFieldString a variable for a neighbor plugin as a string.
func (ss *SecureSite) NeighborPluginFieldString(pluginName string, fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginNeighborfieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.pluginField(pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// NeighborPluginField gets a variable for a neighbor plugin as an interface{}.
func (ss *SecureSite) NeighborPluginField(pluginName string, fieldName string) (interface{}, error) {
	if !ss.Authorized(GrantPluginNeighborfieldRead) {
		return "", ErrAccessDenied
	}

	return ss.pluginField(pluginName, fieldName)
}

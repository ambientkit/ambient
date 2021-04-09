package core

import "fmt"

// PluginNeighborSettingsList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) PluginNeighborSettingsList(pluginName string) ([]Field, error) {
	if !ss.Authorized(GrantPluginNeighborfieldRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	return plugin.Fields(), nil
}

// SetPluginSetting sets a variable for the plugin.
func (ss *SecureSite) SetPluginSetting(settingName string, value string) error {
	if !ss.Authorized(GrantPluginFieldWrite) {
		return ErrAccessDenied
	}

	return ss.pluginsystem.SetSetting(ss.pluginName, settingName, value)
}

// PluginSettingBool returns a plugin setting as a bool.
func (ss *SecureSite) PluginSettingBool(name string) (bool, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return false, ErrAccessDenied
	}

	value, err := ss.settingField(ss.pluginName, name)

	return value == "true", err
}

// PluginSettingString returns a setting for the plugin as a string.
func (ss *SecureSite) PluginSettingString(fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.settingField(ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// PluginSetting returns a setting for the plugin as an interface{}.
func (ss *SecureSite) PluginSetting(fieldName string) (interface{}, error) {
	if !ss.Authorized(GrantPluginFieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.settingField(ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// SetNeighborPluginSetting sets a setting for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginSetting(pluginName string, settingName string, value string) error {
	if !ss.Authorized(GrantPluginNeighborfieldWrite) {
		return ErrAccessDenied
	}

	return ss.pluginsystem.SetSetting(pluginName, settingName, value)
}

// NeighborPluginSettingString returns a setting for a neighbor plugin as a string.
func (ss *SecureSite) NeighborPluginSettingString(pluginName string, fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginNeighborfieldRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.settingField(ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return fmt.Sprint(ival), nil
}

// NeighborPluginSetting returns a setting for a neighbor plugin as an interface{}.
func (ss *SecureSite) NeighborPluginSetting(pluginName string, fieldName string) (interface{}, error) {
	if !ss.Authorized(GrantPluginNeighborfieldRead) {
		return "", ErrAccessDenied
	}

	return ss.settingField(pluginName, fieldName)
}

func (ss *SecureSite) settingField(pluginName string, settingName string) (interface{}, error) {
	raw, err := ss.pluginsystem.Setting(pluginName, settingName)
	if err != nil {
		return "", err
	}

	if raw != nil {
		return raw, nil
	}

	defaultValue, err := ss.pluginsystem.SettingDefault(pluginName, settingName)
	if err != nil {
		return "", err
	}

	return defaultValue, nil
}

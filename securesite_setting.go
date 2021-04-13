package ambient

import "fmt"

// PluginNeighborSettingsList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) PluginNeighborSettingsList(pluginName string) ([]Setting, error) {
	if !ss.Authorized(GrantPluginNeighborSettingRead) {
		return nil, ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, ErrNotFound
	}

	return plugin.Settings(), nil
}

// SetPluginSetting sets a variable for the plugin.
func (ss *SecureSite) SetPluginSetting(settingName string, value string) error {
	if !ss.Authorized(GrantPluginSettingWrite) {
		return ErrAccessDenied
	}

	return ss.pluginsystem.SetSetting(ss.pluginName, settingName, value)
}

// PluginSettingBool returns a plugin setting as a bool.
func (ss *SecureSite) PluginSettingBool(name string) (bool, error) {
	if !ss.Authorized(GrantPluginSettingRead) {
		return false, ErrAccessDenied
	}

	value, err := ss.settingField(ss.pluginName, name)

	return value == "true", err
}

// PluginSettingString returns a setting for the plugin as a string.
func (ss *SecureSite) PluginSettingString(fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginSettingRead) {
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
	if !ss.Authorized(GrantPluginSettingRead) {
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
	if !ss.Authorized(GrantPluginNeighborSettingWrite) {
		return ErrAccessDenied
	}

	settings, err := ss.PluginNeighborSettingsList(pluginName)
	if err != nil {
		return err
	}

	found := false
	for _, setting := range settings {
		if setting.Name == settingName {
			found = true
			break
		}
	}

	if !found {
		ss.log.Debug("setting to set on plugin %v was not specified by the plugin: %v", pluginName, settingName)
		return ErrSettingNotSpecified
	}

	return ss.pluginsystem.SetSetting(pluginName, settingName, value)
}

// NeighborPluginSettingString returns a setting for a neighbor plugin as a string.
func (ss *SecureSite) NeighborPluginSettingString(pluginName string, fieldName string) (string, error) {
	if !ss.Authorized(GrantPluginNeighborSettingRead) {
		return "", ErrAccessDenied
	}

	ival, err := ss.settingField(pluginName, fieldName)
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
	if !ss.Authorized(GrantPluginNeighborSettingRead) {
		return "", ErrAccessDenied
	}

	return ss.settingField(pluginName, fieldName)
}

func (ss *SecureSite) settingField(pluginName string, settingName string) (interface{}, error) {
	raw, err := ss.pluginsystem.Setting(pluginName, settingName)
	if err != nil {
		if err != ErrNotFound {
			return "", err
		}
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

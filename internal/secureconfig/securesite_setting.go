package secureconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

// PluginNeighborSettingsList gets the grants requests for a neighbor plugin.
func (ss *SecureSite) PluginNeighborSettingsList(ctx context.Context, pluginName string) ([]ambient.Setting, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginNeighborSettingRead) {
		return nil, amberror.ErrAccessDenied
	}

	plugin, err := ss.pluginsystem.Plugin(pluginName)
	if err != nil {
		return nil, amberror.ErrNotFound
	}

	return plugin.Settings(ctx), nil
}

// SetPluginSetting sets a variable for the plugin.
func (ss *SecureSite) SetPluginSetting(ctx context.Context, settingName string, value string) error {
	if !ss.Authorized(ctx, ambient.GrantPluginSettingWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetSetting(ss.pluginName, settingName, value)
}

// PluginSettingBool returns a plugin setting as a bool.
func (ss *SecureSite) PluginSettingBool(ctx context.Context, name string) (bool, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginSettingRead) {
		return false, amberror.ErrAccessDenied
	}

	value, err := ss.settingField(ctx, ss.pluginName, name)
	if err != nil {
		return false, err
	}

	if value == nil {
		return false, nil
	}

	return strconv.ParseBool(fmt.Sprint(value))
}

// PluginSettingString returns a setting for the plugin as a string.
func (ss *SecureSite) PluginSettingString(ctx context.Context, fieldName string) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginSettingRead) {
		return "", amberror.ErrAccessDenied
	}

	ival, err := ss.settingField(ctx, ss.pluginName, fieldName)
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
func (ss *SecureSite) PluginSetting(ctx context.Context, fieldName string) (interface{}, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginSettingRead) {
		return "", amberror.ErrAccessDenied
	}

	ival, err := ss.settingField(ctx, ss.pluginName, fieldName)
	if err != nil {
		return "", err
	}

	// Handle nil.
	if ival == nil {
		return "", nil
	}

	return ival, nil
}

// SetNeighborPluginSetting sets a setting for a neighbor plugin.
func (ss *SecureSite) SetNeighborPluginSetting(ctx context.Context, pluginName string, settingName string, value string) error {
	if !ss.Authorized(ctx, ambient.GrantPluginNeighborSettingWrite) {
		return amberror.ErrAccessDenied
	}

	settings, err := ss.PluginNeighborSettingsList(ctx, pluginName)
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
		return amberror.ErrSettingNotSpecified
	}

	return ss.pluginsystem.SetSetting(pluginName, settingName, value)
}

// NeighborPluginSettingString returns a setting for a neighbor plugin as a string.
func (ss *SecureSite) NeighborPluginSettingString(ctx context.Context, pluginName string, fieldName string) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginNeighborSettingRead) {
		return "", amberror.ErrAccessDenied
	}

	ival, err := ss.settingField(ctx, pluginName, fieldName)
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
func (ss *SecureSite) NeighborPluginSetting(ctx context.Context, pluginName string, fieldName string) (interface{}, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginNeighborSettingRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.settingField(ctx, pluginName, fieldName)
}

func (ss *SecureSite) settingField(ctx context.Context, pluginName string, settingName string) (interface{}, error) {
	raw, err := ss.pluginsystem.Setting(pluginName, settingName)
	if err != nil {
		if err != amberror.ErrNotFound {
			return "", err
		}
	}

	if raw != nil {
		return raw, nil
	}

	defaultValue, err := ss.pluginsystem.SettingDefault(ctx, pluginName, settingName)
	if err != nil {
		return "", err
	}

	return defaultValue, nil
}

// PluginTrusted returns whether a plugin is trusted or not.
func (ss *SecureSite) PluginTrusted(ctx context.Context, pluginName string) (bool, error) {
	if !ss.Authorized(ctx, ambient.GrantPluginTrustedRead) {
		return false, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Trusted(pluginName), nil
}

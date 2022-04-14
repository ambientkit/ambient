package ambient

import (
	"context"
	"fmt"
	"regexp"

	"golang.org/x/mod/semver"
)

var (
	// DisallowedPluginNames is a list of disallowed plugin names.
	disallowedPluginNames = map[string]bool{
		"amb":     false,
		"ambient": false,
		"plugin":  false,
	}

	// Only allow starting with a lowercase letter and then containing all
	// lowercase letters and numbers.
	rePluginName = regexp.MustCompile("^[a-z][a-z0-9]*$")
)

// Validate returns an error if the plugin name or version is not valid.
func Validate(ctx context.Context, p PluginCore) error {
	pluginName := p.PluginName(ctx)

	// Don't allow certain plugin names.
	if allowed, ok := disallowedPluginNames[pluginName]; ok && !allowed {
		return fmt.Errorf("plugin name not allowed: %v", pluginName)
	}

	// Don't allow certain plugin name characters.
	if ok := rePluginName.Match([]byte(pluginName)); !ok {
		return fmt.Errorf("plugin name format not allowed: '%v'", pluginName)
	}

	// Ensure version meets https://semver.org/ requirements.
	if ok := semver.IsValid(fmt.Sprintf("v%v", p.PluginVersion())); !ok {
		return fmt.Errorf("plugin (%v) version not in semver format: %v", pluginName, p.PluginVersion())
	}

	return nil
}

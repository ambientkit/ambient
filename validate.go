package ambient

import (
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
func Validate(p PluginCore) error {
	// Don't allow certain plugin names.
	if allowed, ok := disallowedPluginNames[p.PluginName()]; ok && !allowed {
		return fmt.Errorf("plugin name not allowed: %v", p.PluginName())
	}

	// Don't allow certain plugin name characters.
	if ok := rePluginName.Match([]byte(p.PluginName())); !ok {
		return fmt.Errorf("plugin name format not allowed: '%v'", p.PluginName())
	}

	// Ensure version meets https://semver.org/ requirements.
	if ok := semver.IsValid(fmt.Sprintf("v%v", p.PluginVersion())); !ok {
		return fmt.Errorf("plugin (%v) version not in semver format: %v", p.PluginName(), p.PluginVersion())
	}

	return nil
}

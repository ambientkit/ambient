package core

// SettingType is an HTML type of setting.
type SettingType string

const (
	// Input is a standard text field.
	Input SettingType = "input"
	// Textarea is a textarea field.
	Textarea SettingType = "textarea"
	// Checkbox is a checkbox field.
	Checkbox SettingType = "checkbox"
)

// PluginData represents the plugin storage information.
type PluginData struct {
	Enabled  bool           `json:"enabled"`
	Grants   PluginGrants   `json:"grants"`
	Settings PluginSettings `json:"settings"`
}

// PluginGrants represents an unordered map of grants.
type PluginGrants map[Grant]bool

// PluginSettings represents an unordered map of settings.
type PluginSettings map[string]interface{}

// Setting is a plugin settable field.
type Setting struct {
	Name        string             `json:"name"`
	Type        SettingType        `json:"type"`
	Description SettingDescription `json:"description"`
	Hide        bool               `json:"hide"`
	Default     interface{}        `json:"default"`
}

// SettingDescription is a type of description.
type SettingDescription struct {
	Text string `json:"name"`
	URL  string `json:"url"`
}

// PluginRoutes -
type PluginRoutes struct {
	Routes map[string][]Route
}

// Route -
type Route struct {
	Method string
	Path   string
}

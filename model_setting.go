package ambient

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

// SettingType is an HTML type of setting.
type SettingType string

const (
	// Input is a standard text field.
	Input SettingType = "input"
	// InputPassword is a standard password field.
	InputPassword SettingType = "password"
	// Textarea is a textarea field.
	Textarea SettingType = "textarea"
	// Checkbox is a checkbox field.
	Checkbox SettingType = "checkbox"
)

// SettingDescription is a type of description.
type SettingDescription struct {
	Text string `json:"name"`
	URL  string `json:"url"`
}

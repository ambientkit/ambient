package core

// FieldType is a type of field.
type FieldType string

const (
	// Input is a standard text field.
	Input FieldType = "input"
	// Textarea is a textarea field.
	Textarea FieldType = "textarea"
	// Checkbox is a checkbox field.
	Checkbox FieldType = "checkbox"
)

// PluginData represents the plugin storage information.
type PluginData struct {
	Enabled  bool           `json:"enabled"`
	Grants   PluginGrants   `json:"grants"`
	Settings PluginSettings `json:"fields"`
	//Found   bool             `json:"found"`
	//Grants []Grant          `json:"grants"`
	//Fields map[string]Field `json:"fields"`
	//Order  []string         `json:"order"`
}

// Field is a plugin settable field.
type Field struct {
	Name        string           `json:"name"`
	Type        FieldType        `json:"type"`
	Description FieldDescription `json:"description"`
	Default     interface{}      `json:"default"`
}

// FieldList is an array of fields.
type FieldList []Field

// ModelFields returns array of Field.
func (fl FieldList) ModelFields() map[string]Field {
	arr := make(map[string]Field)
	for _, v := range fl {
		arr[v.Name] = v
	}

	return arr
}

// FieldDescription is a type of description.
type FieldDescription struct {
	Text string `json:"name"`
	URL  string `json:"url"`
}

// PluginFields -
// type PluginFields struct {
// 	Fields map[string]interface{} `json:"fields"`
// }

// PluginSettings -
type PluginSettings map[string]interface{}

// PluginGrants -
type PluginGrants map[Grant]bool

// PluginGrants -
// type PluginGrants struct {
// 	Grants map[Grant]bool `json:"grants"`
// }

// PluginRoutes -
type PluginRoutes struct {
	Routes map[string][]Route
}

// Route -
type Route struct {
	Method string
	Path   string
}

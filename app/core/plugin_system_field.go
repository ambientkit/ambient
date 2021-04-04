package core

import "github.com/josephspurrier/ambient/app/model"

const (
	// Input is a standard text field.
	Input FieldType = "input"
	// Textarea is a textarea field.
	Textarea FieldType = "textarea"
)

// FieldList is an array of fields.
type FieldList []Field

// ModelFields returns array of model.Field.
func (fl FieldList) ModelFields() []model.Field {
	arr := make([]model.Field, 0)
	for _, v := range fl {
		arr = append(arr, model.Field{
			Name:        v.Name,
			Type:        model.FieldType(v.Type),
			Description: model.FieldDescription(v.Description),
		})
	}

	return arr
}

// Field is a plugin settable field.
type Field struct {
	Name        string           `json:"name"`
	Type        FieldType        `json:"type"`
	Description FieldDescription `json:"description"`
	Default     string           `json:"default"`
}

// FieldDescription is a type of description.
type FieldDescription struct {
	Text string `json:"name"`
	URL  string `json:"url"`
}

// // Snippet represents an HTML snippet.
// type Snippet struct {
// 	Path     string        `json:"path"`
// 	Location AssetLocation `json:"location"`
// 	Embedded bool          `json:"embedded"`
// 	Replace  []Replace     `json:"replace"`
// 	Auth     AuthType      `json:"auth"`
// 	//Layout   LayoutType    `json:"layout"`
// 	//Attributes []Attribute   `json:"attributes"`
// }

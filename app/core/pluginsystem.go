package core

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/josephspurrier/ambient/app/model"
)

// PluginSystem -
type PluginSystem struct {
	Plugins map[string]IPlugin
	Routes  map[string]IRouteList
}

// NewPluginSystem -
func NewPluginSystem() *PluginSystem {
	return &PluginSystem{
		Plugins: make(map[string]IPlugin),
	}
}

// PluginMeta represents metadata for a plugin that works with the Ambient
// system.
type PluginMeta struct {
	// Name should be globally unique. Only lowercase letters, numbers,
	// and hypens are permitted. Must start with with a letter.
	Name string `json:"name"`
	// Version must follow https://semver.org/.
	Version string `json:"version"`
	// AppVersion is the first compatible version of Ambient that the
	// plugin works with.
	AppVersion string `json:"appversion"`
	// Permissions is which permissions the plugin requests.
	//Permissions []string `json:"permissions"`
}

// AssetLocation is a location where assets can be added.
type AssetLocation string

// AssetType is a type of asset.
type AssetType string

// AuthType is a type of authentication.
type AuthType string

// LayoutType is a type of layout.
type LayoutType string

// FieldType is a type of field.
type FieldType string

const (
	// LocationHead is at the bottom of the HTML <head> section.
	LocationHead AssetLocation = "head"
	// LocationBody is at the bottom of the HTML <body> section.
	LocationBody AssetLocation = "body"
	// LocationMain is at the bottom of the HTML <main> section.
	LocationMain AssetLocation = "main"

	// AssetStylesheet is a stylesheet element.
	AssetStylesheet AssetType = "stylesheet"
	// AssetJavaScript is a javascript element.
	AssetJavaScript AssetType = "javascript"
	// AssetGeneric is a generic element.
	AssetGeneric AssetType = "generic"

	// AllAuth is both anonymous and authenticated users.
	AllAuth AuthType = "all" // Default.
	// AnonymousOnly is only non-authenticated users.
	AnonymousOnly AuthType = "anonymous"
	// AuthenticatedOnly is only authenticated users.
	AuthenticatedOnly AuthType = "authenticated"

	// Page is a page layout.
	Page LayoutType = "page"
	// Post is a post layout.
	Post LayoutType = "post"
	// Dashboard is a dashboard layout.
	Dashboard LayoutType = "dashboard"
	// Bloglist is a bloglist layout.
	Bloglist LayoutType = "bloglist"

	// Input is a standard text field.
	Input FieldType = "input"
	// Textarea is a textarea field.
	Textarea FieldType = "textarea"
)

// FieldDescription is a type of description.
type FieldDescription struct {
	Text string `json:"name"`
	URL  string `json:"url"`
}

// Field is a plugin settable field.
type Field struct {
	Name        string           `json:"name"`
	Type        FieldType        `json:"type"`
	Description FieldDescription `json:"description"`
}

// Snippet represents an HTML snippet.
type Snippet struct {
	Path     string        `json:"path"`
	Location AssetLocation `json:"location"`
	Embedded bool          `json:"embedded"`
	Replace  []Replace     `json:"replace"`
	Auth     AuthType      `json:"auth"`
	//Layout   LayoutType    `json:"layout"`
	//Attributes []Attribute   `json:"attributes"`
}

// Asset represents an HTML asset like a stylesheet or javascript file.
type Asset struct {
	Filetype   AssetType     `json:"filetype"`
	Location   AssetLocation `json:"location"`
	Auth       AuthType      `json:"auth"`
	Attributes []Attribute   `json:"attributes"`
	LayoutOnly []LayoutType  `json:"layout"`

	TagName    string `json:"tagname"`
	ClosingTag bool   `json:"closingtag"`

	External bool      `json:"external"`
	Inline   bool      `json:"inline"`
	Path     string    `json:"path"`
	Replace  []Replace `json:"replace"`
}

// Routable returns true if the file can be served from the embedded filesystem.
func (file Asset) Routable() bool {
	if file.External || file.Inline || file.Filetype == AssetGeneric {
		return false
	}

	return true
}

// SanitizedPath returns an HTML escaped asset path.
func (file Asset) SanitizedPath() string {
	return html.EscapeString(file.Path)
}

// Element returns an HTML element.
func (file *Asset) Element(v IPlugin, assets fs.FS) string {
	// Build the attributes.
	attrs := make([]string, 0)
	for _, attr := range file.Attributes {
		if attr.Value == nil {
			attrs = append(attrs, fmt.Sprintf(`%v`, html.EscapeString(attr.Name)))
		} else {
			attrs = append(attrs, fmt.Sprintf(`%v="%v"`, html.EscapeString(attr.Name), html.EscapeString(fmt.Sprint(attr.Value))))
		}
	}
	attrsJoined := strings.Join(attrs, " ")
	if len(attrsJoined) > 0 {
		// Add a space at the beginning.
		attrsJoined = " " + attrsJoined
	}

	txt := ""
	switch file.Filetype {
	case AssetStylesheet:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				// FIXME: Do something with these.
				fmt.Println(err.Error())
				return ""
			}
			txt = fmt.Sprintf("<style>%v</style>", string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="%v"%v>`, file.SanitizedPath(), attrsJoined)

			} else {
				txt = fmt.Sprintf(`<link rel="stylesheet" href="/plugins/%v/%v?v=%v"%v>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetJavaScript:
		if file.Inline {
			ff, status, err := file.Contents(assets)
			if status != http.StatusOK {
				// FIXME: Do something with these.
				fmt.Println(err.Error())
				return ""
			}
			txt = fmt.Sprintf("<script>%v</script>", string(ff))
		} else {
			if file.External {
				txt = fmt.Sprintf(`<script type="application/javascript" src="%v"%v></script>`, file.SanitizedPath(), attrsJoined)
			} else {
				txt = fmt.Sprintf(`<script type="application/javascript" src="/plugins/%v/%v?v=%v"%v></script>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
			}
		}
	case AssetGeneric:
		if file.ClosingTag {
			txt = fmt.Sprintf(`<%v%v></%v>`, html.EscapeString(file.TagName), attrsJoined, html.EscapeString(file.TagName))
		} else {
			txt = fmt.Sprintf(`<%v%v>`, html.EscapeString(file.TagName), attrsJoined)
		}
	default:
		fmt.Printf("unsupported asset filetype for plugin (%v): %v", v.PluginName(), file.Filetype)
	}

	return txt
}

// Contents returns the text of the file to inline in HTML after doing replace.
func (file *Asset) Contents(assets fs.FS) (ff []byte, status int, err error) {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Open the file.
	f, err := fsys.Open(file.Path)
	if err != nil {
		return nil, http.StatusNotFound, nil
	}
	defer f.Close()

	// Get the contents.
	ff, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Loop over the items to replace.
	for _, rep := range file.Replace {
		ff = bytes.ReplaceAll(ff, []byte(rep.Find), []byte(rep.Replace))
	}

	return ff, http.StatusOK, nil
}

// Attribute represents an HTML attribute.
type Attribute struct {
	Name  string
	Value interface{}
}

// Replace represents text to find and replace.
type Replace struct {
	Find    string
	Replace string
}

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

// IPlugin represents a plugin.
type IPlugin interface {
	PluginName() string
	PluginVersion() string
	Routes()
	Enable(*Toolkit) error
	Disable() error
	Assets() ([]Asset, *embed.FS)
	Fields() []Field
	Middleware() []func(next http.Handler) http.Handler
	SessionManager() (ISession, error)
	//Header() string
	//Body() string
	//SetSettings()
	// Deactivate() error
	// Uninstall() error
}

// IRouteList -
type IRouteList interface {
	Routes() []IRoute
}

// IRoute -
type IRoute interface {
	Method() string
	Path() string
}

// ISettings -
// type ISettings interface {
// 	Add(name string, fieldType string, defaultValue string)
// }

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
	Param(r *http.Request, name string) string
	Error(status int, w http.ResponseWriter, r *http.Request)
}

// IRender represents a template rendered.
type IRender interface {
	PluginDashboard(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, vars map[string]interface{}) (status int, err error)
	PluginPage(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, vars map[string]interface{}) (status int, err error)
	PluginPageContent(w http.ResponseWriter, r *http.Request, content string, vars map[string]interface{}) (status int, err error)
}

// ISession represents a user session.
type ISession interface {
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
	UserAuthenticated(r *http.Request) (bool, error)
	SetUser(r *http.Request, username string)
	RememberMe(r *http.Request, remember bool)
	Logout(r *http.Request)
}

// IPluginLoader -
type IPluginLoader interface {
	LoadSinglePlugin(name string) error
	DisableSinglePlugin(name string) error
}

// ILogger representer the log service for the application.
type ILogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

// Toolkit -
type Toolkit struct {
	Render       IRender
	Router       IRouter
	Security     ISession
	Site         *SecureSite
	PluginLoader IPluginLoader
	Log          ILogger
}

// type Variables struct{}

// func (v *Variables) PostURL() string {
// 	return ""
// }

// Enable -
func (p *PluginMeta) Enable(*Toolkit) error {
	return nil
}

// Disable -
func (p *PluginMeta) Disable() error {
	return nil
}

// Routes -
func (p *PluginMeta) Routes() {}

// Assets -
func (p *PluginMeta) Assets() ([]Asset, *embed.FS) {
	return nil, nil
}

// Fields -
func (p *PluginMeta) Fields() []Field {
	return nil
}

// PluginName -
func (p *PluginMeta) PluginName() string {
	return html.EscapeString(p.Name)
}

// PluginVersion -
func (p *PluginMeta) PluginVersion() string {
	return p.Version
}

// Middleware -
func (p *PluginMeta) Middleware() []func(next http.Handler) http.Handler {
	return nil
}

// SessionManager -
func (p *PluginMeta) SessionManager() (ISession, error) {
	return nil, nil
}

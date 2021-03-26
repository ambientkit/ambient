package core

import (
	"embed"
	"fmt"
	"html"
	"net/http"
	"strings"
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
)

// Snippet represents an HTML snippet.
type Snippet struct {
	Path     string        `json:"path"`
	Location AssetLocation `json:"location"`
	Embedded bool          `json:"embedded"`
	Replace  []Replace     `json:"replace"`
	Auth     AuthType      `json:"auth"`
	Layout   LayoutType    `json:"layout"`
	//Attributes []Attribute   `json:"attributes"`
}

// Asset represents an HTML asset like a stylesheet or javascript file.
type Asset struct {
	Filetype   AssetType     `json:"filetype"`
	Location   AssetLocation `json:"location"`
	Auth       AuthType      `json:"auth"`
	Attributes []Attribute   `json:"attributes"`

	TagName    string `json:"tagname"`
	ClosingTag bool   `json:"closingtag"`

	Embedded bool      `json:"embedded"`
	Path     string    `json:"path"`
	Replace  []Replace `json:"replace"`
}

// SanitizedPath returns an HTML escaped asset path.
func (file Asset) SanitizedPath() string {
	return html.EscapeString(file.Path)
}

// Element returns an HTML element.
func (file *Asset) Element(v IPlugin) string {
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
		if file.Embedded {
			txt = fmt.Sprintf(`<link rel="stylesheet" href="/plugins/%v/%v?v=%v"%v>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
		} else {
			txt = fmt.Sprintf(`<link rel="stylesheet" href="%v"%v>`, file.SanitizedPath(), attrsJoined)
		}
	case AssetJavaScript:
		if file.Embedded {
			txt = fmt.Sprintf(`<script type="application/javascript" src="/plugins/%v/%v?v=%v"%v></script>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
		} else {
			txt = fmt.Sprintf(`<script type="application/javascript" src="%v"%v></script>`, file.SanitizedPath(), attrsJoined)
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

// IPlugin represents a plugin.
type IPlugin interface {
	PluginName() string
	PluginVersion() string
	Routes() error
	Enable(*Toolkit) error
	Assets() ([]Asset, *embed.FS)
	Fields() []string
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
}

// IRender represents a template rendered.
type IRender interface {
	PluginTemplate(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, vars map[string]interface{}) (status int, err error)
}

// ISecurity -
type ISecurity interface {
	SetCSRF(r *http.Request) string
	CSRF(r *http.Request) bool
}

// IPluginLoader -
type IPluginLoader interface {
	LoadSinglePlugin(name string) error
}

// Toolkit -
type Toolkit struct {
	Render       IRender
	Router       IRouter
	Security     ISecurity
	Site         *SecureSite
	PluginLoader IPluginLoader
}

// Enable -
func (p *PluginMeta) Enable(*Toolkit) error {
	return nil
}

// Routes -
func (p *PluginMeta) Routes() error {
	return nil
}

// Assets -
func (p *PluginMeta) Assets() ([]Asset, *embed.FS) {
	return nil, nil
}

// Fields -
func (p *PluginMeta) Fields() []string {
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

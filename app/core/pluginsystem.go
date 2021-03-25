package core

import (
	"embed"
	"html"
	"net/http"
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

const (
	// LocationHead is at the bottom of the HTML <head> section.
	LocationHead AssetLocation = "head"
	// LocationBody is at the bottom of the HTML <body> section.
	LocationBody AssetLocation = "body"

	// FiletypeStylesheet is a stylesheet element.
	FiletypeStylesheet AssetType = "stylesheet"
	// FiletypeJavaScript is a javascript element.
	FiletypeJavaScript AssetType = "javascript"

	// All is both anonymous and authenticated users.
	All AuthType = "all" // Default.
	// AnonymousOnly is only non-authenticated users.
	AnonymousOnly AuthType = "anonymous"
	// AuthenticatedOnly is only authenticated users.
	AuthenticatedOnly AuthType = "authenticated"
)

// Asset represents an HTML asset like a stylesheet or javascript file.
type Asset struct {
	Path       string        `json:"path"`
	Location   AssetLocation `json:"location"`
	Filetype   AssetType     `json:"filetype"`
	Embedded   bool          `json:"embedded"`
	Replace    []Replace     `json:"replace"`
	Auth       AuthType      `json:"auth"`
	Attributes []Attribute   `json:"attributes"`
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

// SanitizedPath returns an HTML escaped asset path.
func (p Asset) SanitizedPath() string {
	return html.EscapeString(p.Path)
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

package ambsystem

import (
	"embed"
	"fmt"
	"html"
	"net/http"
)

// PluginSystem -
type PluginSystem struct {
	Plugins map[string]IPlugin
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

// PluginSettings -
type PluginSettings struct {
	Enabled bool `json:"enabled"`
	Found   bool `json:"found"`
}

type AssetLocation string
type AssetType string

const (
	LocationHeader AssetLocation = "header"
	LocationBody   AssetLocation = "body"

	FiletypeStylesheet AssetType = "stylesheet"
	FiletypeJavaScript AssetType = "javascript"
)

// Asset -
type Asset struct {
	Path     string        `json:"path"`
	Location AssetLocation `json:"location"`
	Filetype AssetType     `json:"filetype"`
	Embedded bool          `json:"embedded"`
}

// SanitizedPath -
func (p Asset) SanitizedPath() string {
	return html.EscapeString(p.Path)
}

// IPlugin represents a plugin.
type IPlugin interface {
	PluginName() string
	PluginVersion() string
	SetPages(IRouter) error
	Assets() ([]Asset, *embed.FS)
	Header() string
	Body() string
	//SetSettings()
	// Deactivate() error
	// Uninstall() error
}

// ISettings -
// type ISettings interface {
// 	Add(name string, fieldType string, defaultValue string)
// }

// IRouter represents a router.
type IRouter interface {
	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
}

// SetPages -
func (p PluginMeta) SetPages(mux IRouter) error {
	fmt.Println("No page to add.")
	return nil
}

// EmbeddedAssets -
func (p PluginMeta) EmbeddedAssets() ([]string, *embed.FS) {
	return nil, nil
}

// Header -
func (p PluginMeta) Header() string {
	return ""
}

// PluginName -
func (p PluginMeta) PluginName() string {
	return html.EscapeString(p.Name)
}

// PluginVersion -
func (p PluginMeta) PluginVersion() string {
	return p.Version
}

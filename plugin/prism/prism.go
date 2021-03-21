package prism

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
)

//go:embed *
var assets embed.FS

// Plugin -
type Plugin struct {
	ambsystem.PluginMeta
}

// New sets up the plugin.
func New() Plugin {
	return Plugin{
		PluginMeta: ambsystem.PluginMeta{
			Name:       "prism",
			Version:    "1.0.0",
			AppVersion: "1.0.0",
		},
	}
}

// Assets -
func (p Plugin) Assets() *embed.FS {
	return &assets
}

// SetPages -
func (p Plugin) SetPages(mux ambsystem.IRouter) error {
	mux.Get("/plugins...", func(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
		// Don't allow directory browsing.
		if strings.HasSuffix(r.URL.Path, "/") {
			return http.StatusNotFound, nil
		}

		// Use the root directory.
		fsys, err := fs.Sub(assets, ".")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Get the requested file name.
		fname := strings.TrimPrefix(r.URL.Path, "/plugins/")

		// Open the file.
		f, err := fsys.Open(fname)
		if err != nil {
			return http.StatusNotFound, nil
		}
		defer f.Close()

		// Get the file time.
		st, err := f.Stat()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.ServeContent(w, r, fname, st.ModTime(), f.(io.ReadSeeker))
		return
	})

	return nil
}

// SetPages -
func (p Plugin) Header() string {
	return `
	<link rel="stylesheet" href="{{"/plugins/css/prism-vsc-dark-plus.css" | AssetStamp}}">
	<style>
	pre[class*="language-"] {
		padding: 0 !important;
	}

	code[class*="language-"] {
		background-color: inherit;
	}
	</style>
	`
}

// Body -
func (pm Plugin) Body() string {
	return `
	<script src="https://unpkg.com/prismjs@1.23.0/components/prism-core.min.js"></script>
	<script src="https://unpkg.com/prismjs@1.23.0/plugins/autoloader/prism-autoloader.min.js"></script>
	`
}

// // SetSettings -
// func (pm Prism) SetSettings(s ambsystem.ISettings) error {
// 	s.Add("name string", fieldType string, defaultValue string)

// }

// Deactivate deactivates the plugin, but leaves the state in the system.
func Deactivate() error {
	return nil
}

// Uninstall removes all plugin state from the system.
func Uninstall() error {
	return nil
}

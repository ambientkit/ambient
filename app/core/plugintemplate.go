package core

import (
	"fmt"
	"html/template"
	"net/http"
)

// InjectPlugins -
func (c *App) InjectPlugins(t *template.Template, r *http.Request) (*template.Template, error) {
	pluginHeader := ""
	pluginBody := ""
	for name, plugin := range c.Storage.Site.PluginSettings {
		if !plugin.Enabled || !plugin.Found {
			continue
		}

		v, found := c.Plugins.Plugins[name]
		if !found {
			fmt.Println("Plugin is missing - should never see this:", name)
			continue
		}

		files, _ := v.Assets()
		if files == nil {
			continue
		}

		_, loggedIn := c.Sess.User(r)
		for _, file := range files {
			// Handle authentication on resources without changing resources.
			if !authAssetAllowed(loggedIn, file) {
				continue
			}

			txt := ""
			switch file.Filetype {
			case FiletypeStylesheet:
				if file.Embedded {
					txt = fmt.Sprintf(`<link rel="stylesheet" href="/plugins/%v/%v?v=%v">`, v.PluginName(), file.SanitizedPath(), v.PluginVersion())
				} else {
					txt = fmt.Sprintf(`<link rel="stylesheet" href="%v">`, file.SanitizedPath())
				}
			case FiletypeJavaScript:
				if file.Embedded {
					txt = fmt.Sprintf(`<script type="application/javascript" src="/plugins/%v/%v?v=%v"></script>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion())
				} else {
					txt = fmt.Sprintf(`<script type="application/javascript" src="%v"></script>`, file.SanitizedPath())
				}
			default:
				fmt.Printf("unsupported asset filetype for plugin (%v): %v", v.PluginName(), file.Filetype)
			}

			switch file.Location {
			case LocationBody:
				pluginBody += txt + "\n    "
			case LocationHeader:
				pluginHeader += txt + "\n    "
			default:
				fmt.Printf("unsupported asset location for plugin (%v): %v", v.PluginName(), file.Filetype)
			}
		}

		//pluginHeader += v.Header()
		//pluginBody += v.Body()
	}

	content := fmt.Sprintf(`{{define "PluginHeaderContent"}}%s{{end}}`, pluginHeader)
	var err error
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginBodyContent"}}%s{{end}}`, pluginBody)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	return t, nil
}

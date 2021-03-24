package core

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strings"
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
			case FiletypeStylesheet:

				if file.Embedded {
					txt = fmt.Sprintf(`<link rel="stylesheet" href="/plugins/%v/%v?v=%v"%v>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
				} else {
					txt = fmt.Sprintf(`<link rel="stylesheet" href="%v"%v>`, file.SanitizedPath(), attrsJoined)
				}
			case FiletypeJavaScript:
				if file.Embedded {
					txt = fmt.Sprintf(`<script type="application/javascript" src="/plugins/%v/%v?v=%v"%v></script>`, v.PluginName(), file.SanitizedPath(), v.PluginVersion(), attrsJoined)
				} else {
					txt = fmt.Sprintf(`<script type="application/javascript" src="%v"%v></script>`, file.SanitizedPath(), attrsJoined)
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

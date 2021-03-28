package core

import (
	"fmt"
	"html/template"
	"net/http"
)

// InjectPlugins will return a template and an error.
func (c *App) InjectPlugins(t *template.Template, r *http.Request, pluginNames []string) (*template.Template, error) {
	pluginHead := ""
	pluginMain := ""
	pluginBody := ""

	// Loop through each of the plugins.
	// Use the plugin names because it's ordered.
	for _, name := range pluginNames {
		plugin, ok := c.Storage.Site.PluginSettings[name]
		if !ok || !plugin.Enabled || !plugin.Found {
			continue
		}

		v, found := c.Plugins.Plugins[name]
		if !found {
			fmt.Println("Plugin is missing - should never see this:", name)
			continue
		}

		files, assets := v.Assets()
		if files == nil {
			continue
		}

		_, loggedIn := c.Sess.User(r)
		for _, file := range files {
			// Handle authentication on resources without changing resources.
			if !authAssetAllowed(loggedIn, file) {
				continue
			}

			// Convert the asset to an element.
			txt := file.Element(v, assets)

			switch file.Location {
			case LocationHead:
				pluginHead += txt + "\n    "
			case LocationBody:
				pluginBody += txt + "\n    "
			case LocationMain:
				pluginMain += txt + "\n    "
			default:
				fmt.Printf("unsupported asset location for plugin (%v): %v", v.PluginName(), file.Filetype)
			}
		}

		//pluginHeader += v.Header()
		//pluginBody += v.Body()
	}

	content := fmt.Sprintf(`{{define "PluginHeadContent"}}%s{{end}}`, pluginHead)
	var err error
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	main := fmt.Sprintf(`{{define "PluginMainContent"}}%s{{end}}`, pluginMain)
	t, err = t.Parse(main)
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

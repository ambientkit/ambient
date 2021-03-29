package core

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
	"github.com/oxtoacart/bpool"
)

// bufpool is used to write out HTML after it's been executed and before it's
// written to the ResponseWriter to catch any partially written templates.
var bufpool *bpool.BufferPool = bpool.NewBufferPool(64)

// InjectPlugins will return a template and an error.
func (c *App) InjectPlugins(t *template.Template, r *http.Request, pluginNames []string, pageURL string) (*template.Template, error) {
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

	// Expose the variables to the plugin templates.
	data := map[string]interface{}{
		"SiteURL": c.Storage.Site.SiteURL(),
		"PageURL": pageURL,
	}

	head, err := templatebuffer.ParseTemplate(pluginHead, data)
	if err != nil {
		return nil, err
	}

	content := fmt.Sprintf(`{{define "PluginHeadContent"}}%s{{end}}`, head)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	main, err := templatebuffer.ParseTemplate(pluginMain, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginMainContent"}}%s{{end}}`, main)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	body, err := templatebuffer.ParseTemplate(pluginBody, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginBodyContent"}}%s{{end}}`, body)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	return t, nil
}

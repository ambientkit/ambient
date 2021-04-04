package core

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// PluginInjector represents a plugin injector.
type PluginInjector struct {
	storage *Storage
	sess    ISession
	plugins *PluginSystem
}

// NewPlugininjector returns a PluginInjector.
func NewPlugininjector(storage *Storage, sess ISession, plugins *PluginSystem) *PluginInjector {
	return &PluginInjector{
		storage: storage,
		sess:    sess,
		plugins: plugins,
	}
}

// Inject will return a template and an error.
func (c *PluginInjector) Inject(t *template.Template, r *http.Request,
	pluginNames []string, layoutType string) (*template.Template, error) {
	pluginHead := ""
	pluginMain := ""
	pluginBody := ""

	// Loop through each of the plugins.
	// Use the plugin names because it's ordered.
	for _, name := range pluginNames {
		plugin, ok := c.storage.Site.PluginSettings[name]
		if !ok || !plugin.Enabled || !plugin.Found {
			continue
		}

		v, found := c.plugins.Plugins[name]
		if !found {
			fmt.Println("Plugin is missing - should never see this:", name)
			continue
		}

		files, assets := v.Assets()
		if files == nil {
			continue
		}

		loggedIn, _ := c.sess.UserAuthenticated(r)
		for _, file := range files {
			// Handle authentication on resources without changing resources.
			if !authAssetAllowed(loggedIn, file) {
				continue
			}

			// Determine if the asset is allowed on the current page type.
			if len(file.LayoutOnly) > 0 {
				allowed := false
				for _, layout := range file.LayoutOnly {
					if string(layout) == layoutType {
						allowed = true
						break
					}
				}
				if !allowed {
					continue
				}
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
		"SiteURL": c.storage.Site.SiteURL(),
		"PageURL": r.URL.Path,
	}

	head, err := templatebuffer.ParseTemplate(pluginHead, nil, data)
	if err != nil {
		return nil, err
	}

	content := fmt.Sprintf(`{{define "PluginHeadContent"}}%s{{end}}`, head)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	main, err := templatebuffer.ParseTemplate(pluginMain, nil, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginMainContent"}}%s{{end}}`, main)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	body, err := templatebuffer.ParseTemplate(pluginBody, nil, data)
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

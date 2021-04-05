package core

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
)

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(t *template.Template, r *http.Request, pluginNames []string, layoutType string) (*template.Template, error)
}

// PluginInjector represents a plugin injector.
type PluginInjector struct {
	storage *Storage
	sess    ISession
	plugins *PluginSystem
	log     ILogger
}

// NewPlugininjector returns a PluginInjector.
func NewPlugininjector(logger ILogger, storage *Storage, sess ISession, plugins *PluginSystem) *PluginInjector {
	return &PluginInjector{
		storage: storage,
		sess:    sess,
		plugins: plugins,
		log:     logger,
	}
}

// Inject will return a template and an error.
func (c *PluginInjector) Inject(t *template.Template, r *http.Request,
	pluginNames []string, layoutType string) (*template.Template, error) {
	pluginHead := ""
	pluginHeader := ""
	pluginMain := ""
	pluginFooter := ""
	pluginBody := ""

	fm := template.FuncMap{}

	// Loop through each of the plugins.
	// Use the plugin names because it's ordered.
	for _, name := range pluginNames {
		plugin, ok := c.storage.Site.PluginSettings[name]
		if !ok || !plugin.Enabled || !plugin.Found {
			continue
		}

		v, found := c.plugins.Plugins[name]
		if !found {
			c.log.Error("plugin injector: plug is missing: %v", name)
			continue
		}

		files, assets, funcMap := v.Assets()
		if files == nil {
			continue
		}

		// If a FuncMap exists, pass request into FuncMap.
		if funcMap != nil {
			afm := funcMap(r)
			for k, v := range afm {
				fm[k] = v
			}
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
			txt := file.Element(c.log, v, assets)

			switch file.Location {
			case LocationHead:
				pluginHead += txt + "\n    "
			case LocationHeader:
				pluginHeader += txt + "\n    "
			case LocationMain:
				pluginMain += txt + "\n    "
			case LocationFooter:
				pluginFooter += txt + "\n    "
			case LocationBody:
				pluginBody += txt + "\n    "
			default:
				c.log.Error("plugin injector: unsupported asset location for plugin (%v): %v", v.PluginName(), file.Filetype)
			}
		}
	}

	// Expose the variables to the plugin templates.
	// TODO: Should we export these variables? Or let plugins set themselves.
	data := map[string]interface{}{
		"SiteURL": c.storage.Site.SiteURL(),
		"PageURL": r.URL.Path,
	}

	head, err := templatebuffer.ParseTemplate(pluginHead, fm, data)
	if err != nil {
		return nil, err
	}

	content := fmt.Sprintf(`{{define "PluginHeadContent"}}%s{{end}}`, head)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	header, err := templatebuffer.ParseTemplate(pluginHeader, fm, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginHeaderContent"}}%s{{end}}`, header)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	main, err := templatebuffer.ParseTemplate(pluginMain, fm, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginMainContent"}}%s{{end}}`, main)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	footer, err := templatebuffer.ParseTemplate(pluginFooter, fm, data)
	if err != nil {
		return nil, err
	}

	content = fmt.Sprintf(`{{define "PluginFooterContent"}}%s{{end}}`, footer)
	t, err = t.Parse(content)
	if err != nil {
		return nil, err
	}

	body, err := templatebuffer.ParseTemplate(pluginBody, fm, data)
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

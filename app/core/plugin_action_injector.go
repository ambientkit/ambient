package core

import (
	"html/template"
	"net/http"
)

// TemplateInjector represents an injector that the template enginer must implement.
type TemplateInjector interface {
	InjectHead(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	InjectHeader(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	InjectMain(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	InjectFooter(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	InjectBody(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
}

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(t *template.Template, r *http.Request, pluginNames []string, layoutType string, vars map[string]interface{}) (*template.Template, error)
}

// PluginInjector represents a plugin injector.
type PluginInjector struct {
	storage *Storage
	sess    ISession
	plugins *PluginSystem
	log     ILogger
	ti      TemplateInjector
}

// NewPlugininjector returns a PluginInjector.
func NewPlugininjector(logger ILogger, ti TemplateInjector, storage *Storage, sess ISession, plugins *PluginSystem) *PluginInjector {
	return &PluginInjector{
		storage: storage,
		sess:    sess,
		plugins: plugins,
		log:     logger,
		ti:      ti,
	}
}

// Inject will return a template and an error.
func (c *PluginInjector) Inject(t *template.Template, r *http.Request, pluginNames []string, layoutType string, vars map[string]interface{}) (*template.Template, error) {
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

	// Add the local variables.
	for k, v := range vars {
		data[k] = v
	}

	// Inject.
	var err error
	t, err = c.ti.InjectHead(t, pluginHead, fm, data)
	if err != nil {
		return nil, err
	}

	t, err = c.ti.InjectHeader(t, pluginHeader, fm, data)
	if err != nil {
		return nil, err
	}

	t, err = c.ti.InjectMain(t, pluginMain, fm, data)
	if err != nil {
		return nil, err
	}

	t, err = c.ti.InjectBody(t, pluginBody, fm, data)
	if err != nil {
		return nil, err
	}

	t, err = c.ti.InjectFooter(t, pluginFooter, fm, data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

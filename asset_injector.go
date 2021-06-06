package ambient

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// AssetInjector represents code that can inject files into a template.
type AssetInjector interface {
	Inject(injector LayoutInjector, t *template.Template, r *http.Request, layoutType LayoutType, vars map[string]interface{}) (*template.Template, error)
	DebugTemplates() bool
}

// LayoutInjector represents an injector that the AssetInjector will call to inject assets in the correct place.
type LayoutInjector interface {
	Head(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Header(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Main(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Footer(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
	Body(t *template.Template, content string, fm template.FuncMap, data map[string]interface{}) (*template.Template, error)
}

// PluginInjector represents a plugin injector.
type PluginInjector struct {
	log            AppLogger
	pluginsystem   *PluginSystem
	sess           AppSession
	debugTemplates bool
}

// NewPlugininjector returns a PluginInjector.
func NewPlugininjector(logger AppLogger, plugins *PluginSystem, sess AppSession, debugTemplates bool) *PluginInjector {
	return &PluginInjector{
		log:            logger,
		pluginsystem:   plugins,
		sess:           sess,
		debugTemplates: debugTemplates,
	}
}

// DebugTemplates returns true if the templates should output debugging information.
func (c *PluginInjector) DebugTemplates() bool {
	return c.debugTemplates
}

// Inject will return a template and an error.
func (c *PluginInjector) Inject(inject LayoutInjector, t *template.Template, r *http.Request, layoutType LayoutType, vars map[string]interface{}) (*template.Template, error) {
	pluginHead := ""
	pluginHeader := ""
	pluginMain := ""
	pluginFooter := ""
	pluginBody := ""

	fm := template.FuncMap{}

	// Loop through each of the plugins.
	// Use the plugin names because it's ordered.
	for _, name := range c.pluginsystem.names {
		plugin, err := c.pluginsystem.PluginData(name)
		if err != nil || !plugin.Enabled {
			continue
		}

		v, err := c.pluginsystem.Plugin(name)
		if err != nil {
			c.log.Error("plugin injector: plugin is missing: %v", name)
			continue
		}

		// If a FuncMap exists, pass request into FuncMap.
		funcMap := v.FuncMap()
		if funcMap != nil {
			// Ensure the plugin has access to write to FuncMap.
			if c.pluginsystem.Authorized(name, GrantSiteFuncMapWrite) {
				afm := funcMap(r)
				for fName, fValue := range afm {
					// Ensure each of the FuncMaps are namespaced.
					if !strings.HasPrefix(fName, v.PluginName()) {
						fm[fmt.Sprintf("%v_%v", v.PluginName(), fName)] = fValue
					} else {
						fm[fName] = fValue
					}
				}
			}
		}

		// Ensure the plugin has access to write to assets.
		files, assets := v.Assets()
		if len(files) > 0 {
			if c.pluginsystem.Authorized(name, GrantSiteAssetWrite) {
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
							if layout == layoutType {
								allowed = true
								break
							}
						}
						if !allowed {
							continue
						}
					}

					// Convert the asset to an element.
					txt := file.Element(c.log, v, assets, c.debugTemplates)

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
		}
	}

	// Add the local variables.
	data := map[string]interface{}{}
	for k, v := range vars {
		data[k] = v
	}

	// Inject into each component.
	var err error
	t, err = inject.Head(t, pluginHead, globalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Header(t, pluginHeader, globalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Main(t, pluginMain, globalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Body(t, pluginBody, globalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Footer(t, pluginFooter, globalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

package injector

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ambientkit/ambient"
)

// PluginInjector represents a plugin injector.
type PluginInjector struct {
	log             ambient.AppLogger
	pluginsystem    ambient.PluginSystem
	sess            ambient.AppSession
	debugTemplates  bool
	escapeTemplates bool
}

// NewPlugininjector returns a PluginInjector.
func NewPlugininjector(logger ambient.AppLogger, plugins ambient.PluginSystem, sess ambient.AppSession, debugTemplates bool, escapeTemplates bool) *PluginInjector {
	return &PluginInjector{
		log:             logger,
		pluginsystem:    plugins,
		sess:            sess,
		debugTemplates:  debugTemplates,
		escapeTemplates: escapeTemplates,
	}
}

// DebugTemplates returns true if the templates should output debugging information.
func (c *PluginInjector) DebugTemplates() bool {
	return c.debugTemplates
}

// EscapeTemplates returns false if template escaping should be disabled.
func (c *PluginInjector) EscapeTemplates() bool {
	return c.escapeTemplates
}

// Inject will return a template and an error.
func (c *PluginInjector) Inject(ctx context.Context, inject ambient.LayoutInjector, t *template.Template,
	r *http.Request, layoutType ambient.LayoutType, vars map[string]interface{}) (*template.Template, error) {
	pluginHead := ""
	pluginHeader := ""
	pluginMain := ""
	pluginFooter := ""
	pluginBody := ""

	fm := template.FuncMap{}

	// Loop through each of the plugins.
	// Use the plugin names because it's ordered.
	for _, name := range c.pluginsystem.Names() {
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
			if c.pluginsystem.Authorized(name, ambient.GrantSiteFuncMapWrite) {
				afm := funcMap(r)
				for fName, fValue := range afm {
					// Ensure each of the FuncMaps are namespaced.
					if !strings.HasPrefix(fName, name) {
						fm[fmt.Sprintf("%v_%v", name, fName)] = fValue
					} else {
						fm[fName] = fValue
					}
				}
			}
		}

		// Ensure the plugin has access to write to assets.
		files, assets := v.Assets()
		if len(files) > 0 {
			if c.pluginsystem.Authorized(name, ambient.GrantSiteAssetWrite) {
				_, err := c.sess.AuthenticatedUser(r)

				for _, file := range files {
					// Handle authentication on resources without changing resources.
					if !ambient.AuthAssetAllowed(err == nil, file) {
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
					txt := file.Element(ctx, c.log, v, assets, c.debugTemplates)

					switch file.Location {
					case ambient.LocationHead:
						if strings.Contains(txt, "charset") {
							// Move charset to the top of the location head.
							// https://webhint.io/docs/user-guide/hints/hint-meta-charset-utf-8/?source=devtools
							pluginHead = txt + "\n    " + pluginHead
						} else {
							// The rest can go after.
							pluginHead += txt + "\n    "
						}
					case ambient.LocationHeader:
						pluginHeader += txt + "\n    "
					case ambient.LocationMain:
						pluginMain += txt + "\n    "
					case ambient.LocationFooter:
						pluginFooter += txt + "\n    "
					case ambient.LocationBody:
						pluginBody += txt + "\n    "
					default:
						c.log.Error("plugin injector: unsupported asset location for plugin (%v): %v", name, file.Filetype)
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
	t, err = inject.Head(t, pluginHead, ambient.GlobalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Header(t, pluginHeader, ambient.GlobalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Main(t, pluginMain, ambient.GlobalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Body(t, pluginBody, ambient.GlobalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	t, err = inject.Footer(t, pluginFooter, ambient.GlobalFuncMap(fm), data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

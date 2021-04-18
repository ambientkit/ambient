# prism

The prism plugin for [Ambient](https://github.com/josephspurrier/ambient) provides syntax highlighting using [Prism](https://prismjs.com/).

**Current version:** 1.0.0

## Example Usage

```go
// Package app represents an Ambient app.
package app

import (
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/prism"
)

// Plugins defines the plugins to use in the application. The order does matter.
var Plugins = func() *ambient.PluginLoader {
	return &ambient.PluginLoader{
		Plugins: []ambient.Plugin{
			// ...
			prism.New(), // Prism CSS for codeblocks.
		},
	}
}
```

## Configuration

### Initialization

The plugin requires these values to be passed in during initialization:

None

### Settings

The plugin allows you to customize these settings:

- **Version**
- **Styles**

### Permissions

The plugin requires these permissions:

- **plugin.setting:read** - Access to add stylesheets and javascript to each page.
- **router.route:write** - Access to create routes for accessing stylesheets.
- **site.asset:write** - Read own plugin settings.

### Environment Variables

The plugin can accept these environment variables:

None
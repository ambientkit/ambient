# prism

The prism plugin for [Ambient](https://github.com/josephspurrier/ambient) provides syntax highlighting using [Prism](https://prismjs.com/).

## Example Usage

```go
// Package app represents an Ambient app.
package app

import (
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/prism"
)

// Plugins defines the plugins to use in the application. The order does matter.
var Plugins = func() ambient.IPluginList {
	return ambient.IPluginList{
		// ...
		prism.New(), // Prism CSS for codeblocks.
	}
}
```

## Configuration

### Required

None

### Optional

None

### Environment Variables

If you want to configure the provider via environment variables, you can use these below.

None
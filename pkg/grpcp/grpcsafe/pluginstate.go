package grpcsafe

import (
	"html/template"
	"sync"

	"github.com/ambientkit/ambient"
	"golang.org/x/net/context"
)

// PluginState contains state used by a plugin.
type PluginState struct {
	contextMap      map[string]context.Context
	contextMapMutex sync.RWMutex

	assetMap      map[string]*AssetContainer
	assetMapMutex sync.RWMutex
}

// AssetContainer contains a FuncMap and a virtual filesystem.
type AssetContainer struct {
	FuncMap template.FuncMap
	FS      ambient.FileSystemReader
}

// NewPluginState returns a thread safe plugin state object.
func NewPluginState() *PluginState {
	return &PluginState{
		contextMap: make(map[string]context.Context),
		assetMap:   make(map[string]*AssetContainer),
	}
}

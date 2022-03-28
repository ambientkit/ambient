package grpcsafe

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

// PluginState contains state used by a plugin.
type PluginState struct {
	contextMap      map[string]context.Context
	contextMapMutex sync.RWMutex

	assetMap      map[string]*AssetContainer
	assetMapMutex sync.RWMutex

	reqMap      map[string]func(http.ResponseWriter, *http.Request) error
	reqMapMutex sync.RWMutex
}

// NewPluginState returns a thread safe plugin state object.
func NewPluginState() *PluginState {
	return &PluginState{
		contextMap: make(map[string]context.Context),
		assetMap:   make(map[string]*AssetContainer),
		reqMap:     make(map[string]func(http.ResponseWriter, *http.Request) error),
	}
}

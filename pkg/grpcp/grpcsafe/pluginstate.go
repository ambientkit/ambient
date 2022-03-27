package grpcsafe

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

// PluginState contains state used by a plugin.
type PluginState struct {
	contextMap       map[string]context.Context
	contextMapMutext sync.RWMutex
}

// NewPluginState returns a thread safe plugin state object.
func NewPluginState() *PluginState {
	return &PluginState{
		contextMap: make(map[string]context.Context),
	}
}

// SaveContext will save the request context.
func (m *PluginState) SaveContext(c context.Context, requestID string) {
	m.contextMapMutext.Lock()
	m.contextMap[requestID] = c
	m.contextMapMutext.Unlock()
}

// Context will return the request context.
func (m *PluginState) Context(requestID string) (context.Context, bool) {
	m.contextMapMutext.RLock()
	c, found := m.contextMap[requestID]
	m.contextMapMutext.RUnlock()
	return c, found
}

// DeleteContextDelayed will delete the request context after 30 seconds.
func (m *PluginState) DeleteContextDelayed(requestID string) {
	go func() {
		<-time.After(30 * time.Second)
		m.contextMapMutext.Lock()
		delete(m.contextMap, requestID)
		m.contextMapMutext.Unlock()
	}()
}

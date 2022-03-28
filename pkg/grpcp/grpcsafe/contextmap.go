package grpcsafe

import (
	"time"

	"golang.org/x/net/context"
)

// SaveContext will save the request context.
func (m *PluginState) SaveContext(c context.Context, requestID string) {
	m.contextMapMutex.Lock()
	m.contextMap[requestID] = c
	m.contextMapMutex.Unlock()
}

// Context will return the request context.
func (m *PluginState) Context(requestID string) (context.Context, bool) {
	m.contextMapMutex.RLock()
	c, found := m.contextMap[requestID]
	m.contextMapMutex.RUnlock()
	return c, found
}

// DeleteContextDelayed will delete the request context after 30 seconds.
func (m *PluginState) DeleteContextDelayed(requestID string) {
	go func() {
		<-time.After(30 * time.Second)
		m.contextMapMutex.Lock()
		delete(m.contextMap, requestID)
		m.contextMapMutex.Unlock()
	}()
}

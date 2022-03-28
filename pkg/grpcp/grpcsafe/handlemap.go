package grpcsafe

import (
	"fmt"
	"net/http"
)

func pathkey(method string, path string) string {
	return fmt.Sprintf("%v %v", method, path)
}

// SaveHandleMap will save the handle map.
func (m *PluginState) SaveHandleMap(c func(http.ResponseWriter, *http.Request) error, method string, path string) {
	m.handleMapMutex.Lock()
	m.handleMap[pathkey(method, path)] = c
	m.handleMapMutex.Unlock()
}

// HandleMap will return the handle map.
func (m *PluginState) HandleMap(method string, path string) (func(http.ResponseWriter, *http.Request) error, bool) {
	m.handleMapMutex.RLock()
	c, found := m.handleMap[pathkey(method, path)]
	m.handleMapMutex.RUnlock()
	return c, found
}

// DeleteHandleMap will delete the handle map.
func (m *PluginState) DeleteHandleMap(requestID string) {
	m.handleMapMutex.Lock()
	delete(m.handleMap, requestID)
	m.handleMapMutex.Unlock()
}

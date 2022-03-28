package grpcsafe

import (
	"html/template"
	"net/http"
	"sync"
)

// HTTPContainer contains the request and response writer for a call.
type HTTPContainer struct {
	Request  *http.Request
	Response http.ResponseWriter
	FuncMap  template.FuncMap
}

// ServerState .
type ServerState struct {
	arr   map[string]*HTTPContainer
	mutex sync.RWMutex
}

// NewServerState .
func NewServerState() *ServerState {
	return &ServerState{
		arr: make(map[string]*HTTPContainer),
	}
}

// Save .
func (m *ServerState) Save(requestID string, c *HTTPContainer) {
	m.mutex.Lock()
	m.arr[requestID] = c
	m.mutex.Unlock()
}

// Load .
func (m *ServerState) Load(requestID string) *HTTPContainer {
	m.mutex.RLock()
	c, found := m.arr[requestID]
	m.mutex.RUnlock()
	if !found {
		return nil
	}
	return c
}

// Delete .
func (m *ServerState) Delete(requestID string) {
	m.mutex.Lock()
	delete(m.arr, requestID)
	m.mutex.Unlock()
}

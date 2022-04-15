package mock

import (
	"context"
	"sync"

	"github.com/ambientkit/ambient"
)

// StoragePlugin represents an Ambient plugin.
type StoragePlugin struct{}

// NewStoragePlugin returns an Ambient plugin that provides memory storage.
func NewStoragePlugin() *StoragePlugin {
	return &StoragePlugin{}
}

// PluginName returns the plugin name.
func (p *StoragePlugin) PluginName(context.Context) string {
	return "mockstorage"
}

// PluginVersion returns the plugin version.
func (p *StoragePlugin) PluginVersion(context.Context) string {
	return "1.0.0"
}

// Storage returns data and session storage.
func (p *StoragePlugin) Storage(logger ambient.Logger) (ambient.DataStorer, ambient.SessionStorer, error) {
	// Use local filesytem for site and session information.
	ds := NewMemoryStore()
	ss := NewMemoryStore()

	return ds, ss, nil
}

// MemoryStore represents a file in memory.
type MemoryStore struct {
	content string
	m       *sync.RWMutex
}

// NewMemoryStore returns a local filesystem object with a file path.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		m: &sync.RWMutex{},
	}
}

// Load returns a file contents from the filesystem.
func (s *MemoryStore) Load() ([]byte, error) {
	s.m.RLock()
	b := []byte(s.content)
	s.m.RUnlock()
	return b, nil
}

// Save writes a file to the filesystem and returns an error if one occurs.
func (s *MemoryStore) Save(b []byte) error {
	s.m.Lock()
	s.content = string(b)
	s.m.Unlock()
	return nil
}

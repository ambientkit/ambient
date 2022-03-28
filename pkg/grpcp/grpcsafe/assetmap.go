package grpcsafe

// SaveAssets will save the request assets.
func (m *PluginState) SaveAssets(c *AssetContainer, requestID string) {
	m.assetMapMutex.Lock()
	m.assetMap[requestID] = c
	m.assetMapMutex.Unlock()
}

// Assets will return the assets.
func (m *PluginState) Assets(requestID string) (*AssetContainer, bool) {
	m.assetMapMutex.RLock()
	c, found := m.assetMap[requestID]
	m.assetMapMutex.RUnlock()
	return c, found
}

// DeleteAssets will delete the assets.
func (m *PluginState) DeleteAssets(requestID string) {
	m.assetMapMutex.Lock()
	delete(m.assetMap, requestID)
	m.assetMapMutex.Unlock()
}

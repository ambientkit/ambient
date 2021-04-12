package core

import (
	"embed"
	"io/fs"
)

// fileExists determines if an embedded file exists.
func fileExists(assets *embed.FS, filename string) bool {
	// Use the root directory.
	fsys, err := fs.Sub(assets, ".")
	if err != nil {
		return false
	}

	// Open the file.
	f, err := fsys.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

// authAssetAllowed return true if the user has access to the asset.
func authAssetAllowed(loggedIn bool, f Asset) bool {
	switch true {
	case f.Auth == AuthOnly && !loggedIn:
		return false
	case f.Auth == AuthOnly && loggedIn:
		return true
	case f.Auth == AuthAnonymousOnly && !loggedIn:
		return true
	case f.Auth == AuthAnonymousOnly && loggedIn:
		return false
	}

	//f.Auth == All:
	return true
}

package secureconfig

import (
	"embed"
	"io/fs"
)

// fileExists determines if an embedded file exists.
func fileExists(assets *embed.FS, filename string) bool {
	// Use the root folder.
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

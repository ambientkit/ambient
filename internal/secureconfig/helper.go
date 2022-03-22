package secureconfig

import (
	"io/fs"

	"github.com/ambientkit/ambient"
)

// fileExists determines if an embedded file exists.
func fileExists(assets ambient.FileSystemReader, filename string) bool {
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

package avfs

import (
	"errors"
	"io/fs"
	"time"
)

// DirEntry is a directory entry.
// type DirEntry interface {
// 	// Name returns the name of the file (or subdirectory) described by the entry.
// 	// This name is only the final element of the path (the base name), not the entire path.
// 	// For example, Name would return "hello.go" not "home/gopher/hello.go".
// 	Name() string
// 	// IsDir reports whether the entry describes a directory.
// 	IsDir() bool
// 	// Type returns the type bits for the entry.
// 	// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
// 	Type() FileMode
// 	// Info returns the FileInfo for the file or subdirectory described by the entry.
// 	// The returned FileInfo may be from the time of the original directory read
// 	// or from the time of the call to Info. If the file has been removed or renamed
// 	// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// 	// If the entry denotes a symbolic link, Info reports the information about the link itself,
// 	// not the link's target.
// 	Info() (FileInfo, error)
// }
type DirEntry struct {
	fi    fs.FileInfo
	files map[string]*File
}

// NewDir handler.
func NewDir(name string) *DirEntry {
	return &DirEntry{
		fi: &FileInfo{
			name:     name,
			isDir:    true,
			contents: []byte{},
			modTime:  time.Now(),
			fileMode: 0,
		},
	}
}

// Stat handler.
func (d *DirEntry) Stat() (fs.FileInfo, error) {
	if d.fi != nil {
		return nil, errors.New("invalid dir")
	}

	return d.fi, nil
}

// Read handler.
func (d *DirEntry) Read([]byte) (int, error) {
	return int(d.fi.Size()), nil
}

// Close handler.
func (d *DirEntry) Close() error {
	return nil
}

// Name handler.
func (d *DirEntry) Name() string {
	return d.fi.Name()
}

// IsDir handler.
func (d *DirEntry) IsDir() bool {
	return true
}

// Type handler.
func (d *DirEntry) Type() fs.FileMode {
	return fs.ModeDir | 0555
}

// Info handler.
func (d *DirEntry) Info() (fs.FileInfo, error) {
	if d.fi != nil {
		return nil, errors.New("invalid directory")
	}

	return d.fi, nil
}

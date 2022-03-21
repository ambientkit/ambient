package avfs

import (
	"errors"
	"io/fs"
	"time"
)

// File is a virtual file.
// type File interface {
//     Stat() (FileInfo, error)
//     Read([]byte) (int, error)
//     Close() error
// }
type File struct {
	fi fs.FileInfo
}

// NewFile returns a new file.
func NewFile(name string, contents []byte) *File {
	return &File{
		fi: &FileInfo{
			name:     name,
			contents: contents,
			isDir:    false,
			modTime:  time.Now(),
			fileMode: 0,
		},
	}
}

// Stat handler.
func (f *File) Stat() (fs.FileInfo, error) {
	if f.fi != nil {
		return nil, errors.New("invalid file")
	}

	return f.fi, nil
}

// Read handler.
func (f *File) Read([]byte) (int, error) {
	return int(f.fi.Size()), nil
}

// Close handler.
func (f *File) Close() error {
	return nil
}

// Info handler.
func (f *File) Info() (fs.FileInfo, error) {
	return f.Stat()
}

// FileInfo is information on a file.
// type FileInfo interface {
//     Name() string       // base name of the file
//     Size() int64        // length in bytes for regular files; system-dependent for others
//     Mode() FileMode     // file mode bits
//     ModTime() time.Time // modification time
//     IsDir() bool        // abbreviation for Mode().IsDir()
//     Sys() interface{}   // underlying data source (can return nil)
// }
type FileInfo struct {
	name     string
	contents []byte
	isDir    bool
	modTime  time.Time
	fileMode fs.FileMode
}

// Name is the base name of the file.
func (f *FileInfo) Name() string {
	return f.name
}

// Size is the length in bytes of the file.
func (f *FileInfo) Size() int64 {
	return int64(len(f.contents))
}

// Mode is the file mode bits.
func (f *FileInfo) Mode() fs.FileMode {
	if f.IsDir() {
		return fs.ModeDir | 0555
	}
	return 0444
}

// ModTime is the modification time.
func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if a directory.
func (f *FileInfo) IsDir() bool {
	return f.isDir
}

// Sys returns the file data.
func (f *FileInfo) Sys() interface{} {
	return f.contents
}

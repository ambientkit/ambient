package avfs

import (
	"errors"
	"io/fs"
)

// FS is a simple virtual file system.
type FS struct {
	files map[string]*File
	dirs  map[string]*DirEntry
}

// NewFS returns a new virtual filesystem.
func NewFS() *FS {
	return &FS{
		files: make(map[string]*File),
		dirs:  make(map[string]*DirEntry),
	}
}

// AddFile to filesystem.
func (f *FS) AddFile(name string, contents []byte) {
	f.files[name] = NewFile(name, contents)
}

// Open opens the named file for reading and returns it as an fs.File.
func (f *FS) Open(name string) (fs.File, error) {
	file, ok := f.files[name]
	if ok {
		return file, nil
	}

	dir, ok := f.dirs[name]
	if ok {
		return dir, nil
		//return &openDir{file, f.readDir(name), 0}, nil
	}
	// if !ok {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	// }
}

// ReadDir reads and returns the entire named directory.
func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	_, ok := f.dirs[name]
	if ok {
		list := make([]fs.DirEntry, 0)
		//for _, range
		return list, nil
		//return &openDir{file, f.readDir(name), 0}, nil
	}

	return nil, &fs.PathError{Op: "read", Path: name, Err: errors.New("not a directory")}
	// file, err := f.Open(name)
	// if err != nil {
	// 	return nil, err
	// }
	// dir, ok := file.(*openDir)
	// if !ok {
	// 	return nil, &fs.PathError{Op: "read", Path: name, Err: errors.New("not a directory")}
	// }
	// list := make([]fs.DirEntry, len(dir.files))
	// for i := range list {
	// 	list[i] = &dir.files[i]
	// }
	//return list, nil
	//return []fs.DirEntry{}, nil
}

// ReadFile reads and returns the content of the named file.
func (f *FS) ReadFile(name string) ([]byte, error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, err
	}
	ofile, ok := file.(*File)
	if !ok {
		return nil, &fs.PathError{Op: "read", Path: name, Err: errors.New("is a directory")}
	}
	return ofile.fi.Sys().([]byte), nil
}

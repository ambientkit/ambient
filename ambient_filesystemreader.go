package ambient

import "io/fs"

// FileSystemReader can be used with embed.FS and avfs.FS.
type FileSystemReader interface {
	Open(name string) (fs.File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
	ReadFile(name string) ([]byte, error)
}

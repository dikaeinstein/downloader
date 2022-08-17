package fsys

import (
	"io/fs"
	"os"
)

// OsFS is an os based filesystem.
type OsFS struct{}

func (OsFS) Create(name string) (fs.File, error) {
	return os.Create(name)
}

func (OsFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// Rename renames file from oldPath to newPath
func (OsFS) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

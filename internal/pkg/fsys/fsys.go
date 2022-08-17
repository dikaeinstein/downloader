package fsys

import (
	"fmt"
	"io/fs"
)

// CreatFS is a filesystem that can create a new file.
type CreatFS interface {
	fs.FS
	Create(name string) (fs.File, error)
}

// RenameFS is a filesystem that can rename a file.
type RenameFS interface {
	fs.FS
	Rename(oldPath, newPath string) error
}

// Create a new file using the given filesystem.
func Create(fsys fs.FS, name string) (fs.File, error) {
	if fsys, ok := fsys.(CreatFS); ok {
		return fsys.Create(name)
	}

	return nil, fmt.Errorf("create %s: operation not supported", name)
}

// Rename a file using the given filesystem.
func Rename(fsys fs.FS, oldPath, newPath string) error {
	if fsys, ok := fsys.(RenameFS); ok {
		return fsys.Rename(oldPath, newPath)
	}

	return fmt.Errorf("rename %s %s: operation not supported", oldPath, newPath)
}

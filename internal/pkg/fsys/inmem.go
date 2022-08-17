package fsys

import (
	"bytes"
	"io/fs"
	"testing/fstest"
)

// InMemFS is an in-memory based filesystem.
type InMemFS struct {
	fstest.MapFS
}

type InMemFile struct {
	file *fstest.MapFile
	buf  *bytes.Buffer
	name string
}

func (f *InMemFile) Close() error { return nil }
func (f *InMemFile) Name() string { return f.name }
func (f *InMemFile) Read(p []byte) (int, error) {
	return f.buf.Read(p)
}
func (f *InMemFile) Stat() (fs.FileInfo, error) { return nil, nil }
func (inMemFile *InMemFile) Write(p []byte) (int, error) {
	n, err := inMemFile.buf.Write(p)
	inMemFile.file.Data = inMemFile.buf.Bytes()
	return n, err
}

// NewInMemFS returns a pointer to a new in-memory FS
func NewInMemFS(mapFs fstest.MapFS) *InMemFS {
	return &InMemFS{mapFs}
}

func (inmem *InMemFS) Open(name string) (fs.File, error) {
	return inmem.MapFS.Open(name)
}

func (inmem *InMemFS) Create(name string) (fs.File, error) {
	buf := new(bytes.Buffer)
	f := &InMemFile{
		buf:  buf,
		file: &fstest.MapFile{Data: buf.Bytes(), Mode: fs.FileMode(0o666)},
		name: name,
	}

	inmem.MapFS[name] = f.file

	return f, nil
}

func (inmem *InMemFS) Rename(oldPath, newPath string) error {
	f := inmem.MapFS[oldPath]
	inmem.MapFS[newPath] = f
	delete(inmem.MapFS, oldPath)

	return nil
}

// Content returns the data of the named file as a buffer of bytes.
func (inmem *InMemFS) Content(name string) *bytes.Buffer {
	return bytes.NewBuffer(inmem.MapFS[name].Data)
}

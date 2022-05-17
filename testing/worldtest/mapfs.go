package worldtest

import (
	"io/fs"
	"testing/fstest"
	"time"
)

type MapFS struct {
	fstest.MapFS
	Cwd struct {
		string
		error
	}
}

func NewMapFS(m map[string][]byte) *MapFS {
	mapFS := fstest.MapFS{}
	for k, v := range m {
		mapFS[k] = &fstest.MapFile{
			Data: v,
		}
	}
	return &MapFS{
		MapFS: mapFS,
	}
}

func (m *MapFS) Getwd() (string, error) {
	return m.Cwd.string, m.Cwd.error
}

func (m *MapFS) PathSeparator() rune {
	return '/'
}

func (m *MapFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	m.MapFS[name] = &fstest.MapFile{
		Data:    data,
		Mode:    perm,
		ModTime: time.Now(),
		Sys:     nil,
	}
	return nil
}

func (m *MapFS) ReadFile(name string) ([]byte, error) {
	return m.MapFS.ReadFile(name)
}

func (m *MapFS) Remove(name string) error {
	delete(m.MapFS, name)
	return nil
}

func (m *MapFS) Stat(name string) (fs.FileInfo, error) {
	return m.MapFS.Stat(name)
}

func (m *MapFS) Glob(name string) ([]string, error) {
	return m.MapFS.Glob(name)
}

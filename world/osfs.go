package world

import (
	"io/fs"
	"os"
	"path/filepath"
)

type OSFS struct{}

func (OSFS) Getwd() (string, error) {
	return os.Getwd()
}

func (OSFS) PathSeparator() rune {
	return os.PathSeparator
}

func (OSFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, os.FileMode(perm))
}

func (OSFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OSFS) Remove(name string) error {
	return os.Remove(name)
}

func (OSFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (OSFS) Glob(name string) ([]string, error) {
	return filepath.Glob(name)
}

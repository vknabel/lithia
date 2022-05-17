package world

import (
	"io"
	"io/fs"
)

type World struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	FS  WorldFS
	Env WorldEnv
}

type WorldFS interface {
	Getwd() (string, error)
	PathSeparator() rune
	WriteFile(name string, data []byte, perm fs.FileMode) error
	ReadFile(name string) ([]byte, error)
	Remove(name string) error
	Stat(name string) (fs.FileInfo, error)
	Glob(pattern string) (matches []string, err error)
}

type WorldEnv interface {
	Exit(code int)
	LookupEnv(key string) (string, bool)
}

var Current World = New()

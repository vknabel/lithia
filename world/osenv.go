package world

import "os"

type OSEnv struct{}

func (OSEnv) Exit(code int) {
	os.Exit(code)
}

func (OSEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (OSEnv) Environ() []string {
	return os.Environ()
}

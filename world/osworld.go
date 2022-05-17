package world

import "os"

func New() World {
	return World{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		FS:     OSFS{},
		Env:    OSEnv{},
	}
}

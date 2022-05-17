package worldtest

import (
	"os"

	"github.com/vknabel/lithia/world"
)

func NewTestWorld(env map[string]string, files map[string][]byte) world.World {
	return world.World{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    NewMapEnv(env),
		FS:     NewMapFS(files),
	}
}

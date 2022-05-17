package main

import (
	"fmt"

	"github.com/vknabel/lithia/app/lithia/cmd"
	"github.com/vknabel/lithia/world"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/vknabel/lithia/app/lithia/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

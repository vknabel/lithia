package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/vknabel/lithia"
	"github.com/vknabel/lithia/world"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [script]",
	Short: "Runs a Lithia script",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runFile(args[0], args)
	},
}

func runFile(fileName string, args []string) {
	world.Current.Args = args
	scriptData, err := world.Current.FS.ReadFile(fileName)
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}

	inter, _ := lithia.NewDefaultInterpreter(path.Dir(fileName))
	script := string(scriptData) + "\n"
	_, err = inter.Interpret(fileName, script)
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
}

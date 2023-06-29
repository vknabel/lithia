package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/vknabel/lithia"
	"github.com/vknabel/lithia/potfile"
	"github.com/vknabel/lithia/world"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:                "run [script]",
	Short:              "Runs a Lithia script",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		world.Current.Args = args

		firstArg := args[0]
		potfileState, err := potfile.ForReferenceFile(firstArg)
		if err != nil {
			fmt.Fprint(world.Current.Stderr, err)
			world.Current.Env.Exit(1)
		}

		if potCmd, ok := potfileState.Cmds[firstArg]; ok {
			potCmd.RunCmd(args[1:])
			return
		}

		runFile(firstArg, args)
	},
}

func runFile(fileName string, args []string) {
	scriptData, err := world.Current.FS.ReadFile(fileName)
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
	inter := lithia.NewDefaultInterpreter(path.Dir(fileName))
	script := string(scriptData)
	_, err = inter.Interpret(fileName, script)

	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
}

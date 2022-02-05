package cmd

import (
	"fmt"
	"os"
	"path"

	cobra "github.com/muesli/coral"
	"github.com/vknabel/lithia/runtime"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [script]",
	Short: "Runs a Lithia script",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runFile(args[0])
	},
}

func runFile(fileName string) {
	scriptData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	inter := runtime.NewInterpreter(path.Dir(fileName))
	script := string(scriptData) + "\n"
	_, err = inter.Interpret(fileName, script)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

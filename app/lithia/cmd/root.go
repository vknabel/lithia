package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "lithia",
	Short: "Lithia programming language",
	Long: "Lithia is an experimental functional programming language " +
		"with an implicit but strong and dynamic type system.\n" +
		"It is designed around a few core concepts in mind " +
		"all language features contribute to.\n" +
		"\n" +
		"Lean more at https://github.com/vknabel/lithia",
	Version: "0.0.12-next",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			runFile(args[0])
		} else {
			runPrompt()
		}
	},
}

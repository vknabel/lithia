package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vknabel/lithia/info"
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
	Version: fmt.Sprintf("%s\ncommit: %s\nbuilt by: %s\nbuilt at: %s", info.Version, info.Commit, info.Date, info.BuiltBy),
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			runFile(args[0], args)
		} else {
			runPrompt()
		}
	},
}

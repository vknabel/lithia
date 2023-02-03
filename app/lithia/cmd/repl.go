package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/vknabel/lithia"
	"github.com/vknabel/lithia/reporting"
	"github.com/vknabel/lithia/world"
)

func init() {
	rootCmd.AddCommand(replCmd)
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Runs interactive Lithia REPL.",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runPrompt()
	},
}

func runPrompt() {
	importRoot, err := world.Current.FS.Getwd()
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
	reader := bufio.NewReader(world.Current.Stdin)
	inter, ctx := lithia.NewDefaultInterpreter(importRoot)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		value, err := inter.InterpretEmbed("prompt", line)
		if err != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		if value != nil {
			fmt.Println("- ", value)
		}
	}
}

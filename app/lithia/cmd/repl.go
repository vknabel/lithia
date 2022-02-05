package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	cobra "github.com/muesli/coral"
	"github.com/vknabel/lithia/reporting"
	"github.com/vknabel/lithia/runtime"
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
	importRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	inter := runtime.NewInterpreter(importRoot)
	for {
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

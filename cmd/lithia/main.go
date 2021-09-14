package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/vknabel/go-lithia/interpreter"
	"github.com/vknabel/go-lithia/reporting"
)

func main() {
	args := os.Args[1:]
	var err error
	if len(args) > 1 {
		fmt.Fprint(os.Stderr, "Usage: lithia <script>")
		os.Exit(64)
	} else if len(args) == 1 {
		script := args[0]
		err = runFile(script)
	} else {
		err = runPrompt()
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(65)
	}
}

func runFile(fileName string) error {
	scriptData, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	inter := interpreter.NewInterpreter()
	script := string(scriptData)
	_, err = inter.Interpret(fileName, script)
	return err
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	inter := interpreter.NewInterpreter()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		lazyValue, err := inter.Interpret("prompt", line)
		if err != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		value, err := lazyValue.Evaluate()
		if err != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		if value != nil {
			fmt.Println(value)
		}
	}
}

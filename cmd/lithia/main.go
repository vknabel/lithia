package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vknabel/go-lithia/reporting"
	"github.com/vknabel/go-lithia/runtime"
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
	inter := runtime.NewInterpreter(path.Dir(fileName))
	script := string(scriptData) + "\n"
	_, ierr := inter.Interpret(fileName, script)
	return ierr
}

func runPrompt() error {
	importRoot, err := os.Getwd()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	inter := runtime.NewInterpreter(importRoot)
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
		value, ierr := inter.InterpretEmbed("prompt", line)
		if ierr != nil {
			reporting.ReportErrorOrPanic(err)
			continue
		}
		if value != nil {
			fmt.Println(value)
		}
	}
}

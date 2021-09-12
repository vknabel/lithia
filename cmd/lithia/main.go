package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/interpreter"
	"github.com/vknabel/go-lithia/parser"
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

func runFile(filename string) error {
	scriptData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	script := string(scriptData)
	runScript(script, interpreter.NewInterpreter())
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	interpreter := interpreter.NewInterpreter()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			reporting.ReportError(1, err.Error())
			continue
		}
		lazyValue, err := runScript(line, interpreter)
		if err != nil {
			reporting.ReportError(1, err.Error())
			continue
		}
		value, err := lazyValue.Evaluate()
		if err != nil {
			reporting.ReportError(1, err.Error())
			continue
		}
		if value != nil {
			fmt.Println(value)
		}
	}
}

func runScript(script string, interpreter *interpreter.Interpreter) (*interpreter.LazyRuntimeValue, error) {
	tree, err := parse(script)
	if err != nil {
		return nil, err
	}
	value, err := interpreter.Interpret(tree, []byte(script))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func parse(script string) (*sitter.Tree, error) {
	parser := parser.NewParser()
	tree, error := parser.Parse(script)
	return tree, error
}

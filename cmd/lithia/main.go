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
		} else if err != nil {
			reporting.ReportError(1, err.Error())
		}
		lazyValue := runScript(line, interpreter)
		value, err := lazyValue.Evaluate()
		if err != nil {
			reporting.ReportError(1, err.Error())
		}
		if value != nil {
			fmt.Println(value)
		}
	}
}

func runScript(script string, interpreter *interpreter.Interpreter) *interpreter.LazyRuntimeValue {
	tree, err := parse(script)
	if err != nil {
		reporting.ReportError(0, err.Error())
		return nil
	}
	value, err := interpreter.Interpret(tree, []byte(script))
	if err != nil {
		reporting.ReportError(0, err.Error())
		return nil
	}
	return value
}

func parse(script string) (*sitter.Tree, error) {
	parser := parser.NewParser()
	tree, error := parser.Parse(script)
	return tree, error
}

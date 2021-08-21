package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vknabel/go-lithia/reporting"
	"github.com/vknabel/go-lithia/scanner"
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
	err = run(script)
	return err
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		} else if err != nil {
			reporting.ReportError(1, err.Error())
		}
		run(line)
	}
}

func run(script string) error {
	reader := strings.NewReader(script)
	scanner := scanner.NewScanner(reader)
	tokens, errs := scanner.ScanTokens()

	if len(errs) > 0 {
		for _, err := range errs {
			reporting.ReportError(1, err.Error())
		}
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

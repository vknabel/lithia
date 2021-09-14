package interpreter

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func FormatNodeErrorMessage(kind string, message string, fileName string, node *sitter.Node, source string) string {
	lines := strings.Split(source, "\n")
	currentError := fmt.Sprintf("%s:%d:%d: %s\n", fileName, node.StartPoint().Row+1, node.StartPoint().Column+1, kind)
	minErrorPrefix := strings.Repeat(" ", 4)
	errorPrefix := strings.Repeat(" ", int(node.StartPoint().Column))
	currentError += fmt.Sprintln(minErrorPrefix + lines[node.StartPoint().Row])
	currentError += fmt.Sprintln(minErrorPrefix + errorPrefix + "^")
	currentError += fmt.Sprintln(minErrorPrefix + message)
	return currentError
}

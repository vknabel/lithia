package interpreter

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

type SyntaxError struct {
	Message string
	Node    *sitter.Node
	Source  []byte
	File    string
}

func (ex *EvaluationContext) SyntaxErrorf(format string, args ...interface{}) SyntaxError {
	return SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Node:    ex.node,
		Source:  ex.source,
		File:    ex.file,
	}
}

func (e SyntaxError) Error() string {
	return FormatNodeErrorMessage("syntax error", e.Message, e.File, e.Node, string(e.Source))
}

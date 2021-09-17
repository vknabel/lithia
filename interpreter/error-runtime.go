package interpreter

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

type RuntimeError struct {
	Message string
	Node    *sitter.Node
	Source  []byte
	File    string
}

func (ex *EvaluationContext) RuntimeErrorf(format string, args ...interface{}) RuntimeError {
	return RuntimeError{
		Message: fmt.Sprintf(format, args...),
		Node:    ex.node,
		Source:  ex.source,
		File:    ex.file,
	}
}

func (e RuntimeError) Error() string {
	return FormatNodeErrorMessage("runtime error", e.Message, e.File, e.Node, string(e.Source))
}

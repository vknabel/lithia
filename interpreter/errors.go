package interpreter

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

type SyntaxError struct {
	Message    string
	Node       *sitter.Node
	Source     []byte
	ModuleFile ModuleFile
}

func (ex *ExecutionContext) SyntaxErrorf(format string, args ...interface{}) SyntaxError {
	return SyntaxError{
		Message:    fmt.Sprintf(format, args...),
		Node:       ex.node,
		Source:     ex.source,
		ModuleFile: ex.moduleFile,
	}
}

func (e SyntaxError) Error() string {
	return FormatNodeErrorMessage("syntax error", e.Message, e.ModuleFile.name, e.Node, string(e.Source))
}

type RuntimeError struct {
	Message    string
	Node       *sitter.Node
	Source     []byte
	ModuleFile ModuleFile
}

func (ex *ExecutionContext) RuntimeErrorf(format string, args ...interface{}) RuntimeError {
	return RuntimeError{
		Message:    fmt.Sprintf(format, args...),
		Node:       ex.node,
		Source:     ex.source,
		ModuleFile: ex.moduleFile,
	}
}

func (e RuntimeError) Error() string {
	return FormatNodeErrorMessage("runtime error", e.Message, e.ModuleFile.name, e.Node, string(e.Source))
}

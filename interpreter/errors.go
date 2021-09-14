package interpreter

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/reporting"
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
	return fmt.Sprintf("%s syntax error: %s", e.SourceLocation(), e.Message)
}

func (e SyntaxError) SourceLocation() string {
	var _ reporting.LocatableError = e
	return fmt.Sprintf("%s:%d:%d", e.ModuleFile.name, e.Node.StartPoint().Row+1, e.Node.StartPoint().Column)
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
	return fmt.Sprintf("%s runtime error: %s", e.SourceLocation(), e.Message)
}

func (e RuntimeError) SourceLocation() string {
	var _ reporting.LocatableError = e
	return fmt.Sprintf("%s:%d:%d", e.ModuleFile.name, e.Node.StartPoint().Row+1, e.Node.StartPoint().Column)
}

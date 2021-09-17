package interpreter

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type LocatableError interface {
	error
	LocatableError() LocatableError
}

type SourceLocatableError struct {
	Kind           string
	Message        string
	SourceLocation SourceLocation
}

type SourceLocation struct {
	FileName string
	Source   string
	Node     *sitter.Node
}

func (err SourceLocatableError) Error() string {
	location := err.SourceLocation
	lines := strings.Split(location.Source, "\n")
	currentError := fmt.Sprintf("%s:%d:%d: %s\n", location.FileName, location.Node.StartPoint().Row+1, location.Node.StartPoint().Column+1, err.Kind)
	minErrorPrefix := strings.Repeat(" ", 4)
	errorPrefix := strings.Repeat(" ", int(location.Node.StartPoint().Column))
	currentError += fmt.Sprintln(minErrorPrefix + lines[location.Node.StartPoint().Row])
	currentError += fmt.Sprintln(minErrorPrefix + errorPrefix + "^")
	currentError += fmt.Sprintln(minErrorPrefix + err.Message)
	return currentError
}

func (err SourceLocatableError) LocatableError() LocatableError {
	return err
}

func NewLocatableError(kind string, message string, fileName string, source string, node *sitter.Node) SourceLocatableError {
	return SourceLocatableError{
		Kind:    kind,
		Message: message,
		SourceLocation: SourceLocation{
			Node:     node,
			Source:   source,
			FileName: fileName,
		},
	}
}

func (ex *EvaluationContext) LocatableErrorOrConvert(err error) LocatableError {
	if err == nil {
		return nil
	}
	if locatableError, ok := err.(LocatableError); ok {
		return locatableError
	}
	return NewLocatableError("error", err.Error(), ex.file, string(ex.source), ex.node)
}

func (ex *EvaluationContext) LocatableErrorf(kind string, format string, args ...interface{}) LocatableError {
	return NewLocatableError(kind, fmt.Sprintf(format, args...), ex.file, string(ex.source), ex.node)
}

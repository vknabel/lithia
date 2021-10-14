package parser

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type SyntaxError struct {
	Kind           string
	Message        string
	SourceLocation SourceLocation
}

type SourceLocation struct {
	FileName string
	Source   string
	Node     *sitter.Node
}

func (err SyntaxError) Error() string {
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

func NewSyntaxError(kind string, message string, fileName string, source string, node *sitter.Node) SyntaxError {
	return SyntaxError{
		Kind:    kind,
		Message: message,
		SourceLocation: SourceLocation{
			Node:     node,
			Source:   source,
			FileName: fileName,
		},
	}
}

func (ex *FileParser) SyntaxErrorOrConvert(err error) *SyntaxError {
	if err == nil {
		return nil
	}
	if locatableError, ok := err.(SyntaxError); ok {
		return &locatableError
	}
	if parsingError, ok := err.(SyntaxParsingError); ok {
		syntaxErr := NewSyntaxError("syntax", parsingError.Error(), ex.File, string(ex.Source), ex.Node)
		return &syntaxErr
	}
	syntaxErr := NewSyntaxError("error", err.Error(), ex.File, string(ex.Source), ex.Node)
	return &syntaxErr
}

func (ex *FileParser) LocatableErrorf(kind string, format string, args ...interface{}) SyntaxError {
	return NewSyntaxError(kind, fmt.Sprintf(format, args...), ex.File, string(ex.Source), ex.Node)
}

func (ex *FileParser) SyntaxErrorf(format string, args ...interface{}) SyntaxError {
	return ex.LocatableErrorf("syntax error", format, args...)
}

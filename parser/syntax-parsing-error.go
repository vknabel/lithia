package parser

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type SyntaxParsingError struct {
	fileName string
	source   string
	tree     *sitter.Tree
}

func MakeSyntaxParsingError(fileName string, source string, tree *sitter.Tree) *SyntaxParsingError {
	return &SyntaxParsingError{
		fileName: fileName,
		source:   source,
		tree:     tree,
	}
}

func (err SyntaxParsingError) SyntaxErrors() []SyntaxError {
	partialErrors := []SyntaxError{}

	for _, errorNode := range err.ErrorNodes(err.tree.RootNode()) {
		var message string
		if errorNode.ChildCount() > 0 {
			message = fmt.Sprintln(errorNode.Child(0).String())
		} else {
			message = ""
		}
		currentError := NewSyntaxError("syntax error", message, err.fileName, err.source, errorNode)
		partialErrors = append(partialErrors, currentError)
	}
	return partialErrors
}

func (err SyntaxParsingError) Error() string {
	partials := []string{}
	for _, partialError := range err.SyntaxErrors() {
		partials = append(partials, partialError.Error())
	}
	if len(partials) > 0 {
		return strings.Join(partials, "\n\n")
	} else {
		return fmt.Sprintf("%s: %s\n\n", err.fileName, err.tree.RootNode().String())
	}
}

func (e SyntaxParsingError) ErrorNodes(node *sitter.Node) []*sitter.Node {
	if node.Type() == TYPE_NODE_ERROR {
		return []*sitter.Node{node}
	}

	partial := []*sitter.Node{}
	if node.IsMissing() {
		partial = append(partial, node)
	}
	if node.IsNull() {
		partial = append(partial, node)
	}
	if node.Type() == TYPE_NODE_UNEXPECTED {
		partial = append(partial, node)
	}
	if node.Type() == TYPE_NODE_MISSING {
		partial = append(partial, node)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		partial = append(partial, e.ErrorNodes(child)...)
	}
	return partial
}

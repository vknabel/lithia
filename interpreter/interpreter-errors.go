package interpreter

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type SyntaxParsingError struct {
	fileName string
	source   string
	tree     *sitter.Tree
}

func (inter *Interpreter) SyntaxParsingError(fileName string, source string, tree *sitter.Tree) SyntaxParsingError {
	return SyntaxParsingError{
		fileName: fileName,
		source:   source,
		tree:     tree,
	}
}

func (e SyntaxParsingError) Error() string {
	partialErrors := []string{}

	for _, errorNode := range e.ErrorNodes(e.tree.RootNode()) {
		var message string
		if errorNode.ChildCount() > 0 {
			message = fmt.Sprintln(errorNode.Child(0).String())
		} else {
			message = ""
		}
		currentError := FormatNodeErrorMessage("syntax error", message, e.fileName, errorNode, e.source)
		partialErrors = append(partialErrors, currentError)
	}
	return strings.Join(partialErrors, "\n\n")
}

func (e SyntaxParsingError) ErrorNodes(node *sitter.Node) []*sitter.Node {
	if node.Type() == parser.TYPE_NODE_ERROR {
		return []*sitter.Node{node}
	}

	partial := []*sitter.Node{}
	if node.IsMissing() {
		partial = append(partial, node)
	}
	if node.IsNull() {
		partial = append(partial, node)
	}
	if node.Type() == parser.TYPE_NODE_UNEXPECTED {
		partial = append(partial, node)
	}
	if node.Type() == parser.TYPE_NODE_MISSING {
		partial = append(partial, node)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		partial = append(partial, e.ErrorNodes(child)...)
	}
	return partial
}

package parser

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	syntax "github.com/vknabel/tree-sitter-lithia"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) Parse(contents string) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(syntax.GetLanguage())

	input := []byte(contents)
	tree := parser.Parse(nil, input)

	if tree.RootNode().HasError() {
		return tree, fmt.Errorf("error parsing tree: %s", tree.RootNode())
	}
	return tree, nil
}

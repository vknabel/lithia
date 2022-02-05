package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/lithia/ast"
	syntax "github.com/vknabel/tree-sitter-lithia"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) Parse(moduleName ast.ModuleName, file string, contents string) (*FileParser, []SyntaxError) {
	parser := sitter.NewParser()
	parser.SetLanguage(syntax.GetLanguage())

	input := []byte(contents)
	tree := parser.Parse(nil, input)

	fileParser := NewFileParser(moduleName, file, tree.RootNode(), tree, input)
	if tree.RootNode().HasError() {
		return fileParser, MakeSyntaxParsingError(file, contents, tree).SyntaxErrors()
	}
	return fileParser, nil
}

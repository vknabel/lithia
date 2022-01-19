package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/ast"
)

type FileParser struct {
	ModuleName    ast.ModuleName
	File          string
	FunctionCount int

	Node   *sitter.Node
	Tree   *sitter.Tree
	Source []byte

	Comments []string
}

func NewFileParser(moduleName ast.ModuleName, file string, node *sitter.Node, tree *sitter.Tree, source []byte) *FileParser {
	if node == nil {
		panic("node is nil")
	}
	return &FileParser{
		ModuleName: moduleName,
		File:       file,
		Node:       node,
		Tree:       tree,
		Source:     source,
	}
}

func (fp *FileParser) ConsumeDocs() *ast.Docs {
	docs := ast.MakeDocs(fp.Comments)
	fp.Comments = []string{}
	return docs
}

func (fp *FileParser) ChildParser(node *sitter.Node) *FileParser {
	if node == nil {
		panic("node is nil")
	}
	return &FileParser{
		ModuleName:    fp.ModuleName,
		File:          fp.File,
		FunctionCount: 0,
		Node:          node,
		Source:        fp.Source,
		Comments:      []string{},
	}
}

func (fp *FileParser) ChildParserConsumingComments(node *sitter.Node) *FileParser {
	if node == nil {
		panic("node is nil")
	}
	comments := fp.Comments
	fp.Comments = []string{}
	return &FileParser{
		ModuleName:    fp.ModuleName,
		File:          fp.File,
		FunctionCount: 0,
		Node:          node,
		Source:        fp.Source,
		Comments:      comments,
	}
}

func (fp *FileParser) AstSource() *ast.Source {
	return ast.MakeSource(
		fp.ModuleName,
		fp.File,
		ast.MakePosition(int(fp.Node.StartPoint().Row), int(fp.Node.StartPoint().Column)),
		ast.MakePosition(int(fp.Node.EndPoint().Row), int(fp.Node.EndPoint().Column)),
	)
}

func (fp *FileParser) AssertNodeType(nodeType string) []SyntaxError {
	if fp.Node.Type() != nodeType {
		return []SyntaxError{fp.SyntaxErrorf("unexpected %q, expected %q", fp.Node.Type(), nodeType)}
	}
	return nil
}

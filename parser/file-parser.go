package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/lithia/ast"
)

type FileParser struct {
	ModuleName    ast.ModuleName
	File          string
	functionCount *int

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
		ModuleName:    moduleName,
		File:          file,
		functionCount: new(int),
		Node:          node,
		Tree:          tree,
		Source:        source,
	}
}

func (fp *FileParser) ConsumeDocs() *ast.Docs {
	docs := ast.MakeDocs(fp.Comments)
	fp.Comments = []string{}
	return docs
}

func (fp *FileParser) NewScopeChildParser(node *sitter.Node) *FileParser {
	if node == nil {
		panic("node is nil")
	}
	return &FileParser{
		ModuleName:    fp.ModuleName,
		File:          fp.File,
		functionCount: new(int),
		Node:          node,
		Source:        fp.Source,
		Comments:      []string{},
	}
}

func (fp *FileParser) SameScopeChildParser(node *sitter.Node) *FileParser {
	if node == nil {
		panic("node is nil")
	}
	return &FileParser{
		ModuleName:    fp.ModuleName,
		File:          fp.File,
		functionCount: fp.functionCount,
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
		functionCount: fp.functionCount,
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

func (fp *FileParser) CountFunction() int {
	count := fp.functionCount
	*fp.functionCount = *fp.functionCount + 1
	return *count
}

func (fp *FileParser) addAllChildComments() {
	caseCount := int(fp.Node.NamedChildCount())
	for i := 0; i < caseCount; i++ {
		childNode := fp.Node.NamedChild(i)
		if childNode.Type() == TYPE_NODE_COMMENT {
			fp.Comments = append(fp.Comments, childNode.Content(fp.Source))
		}
	}
}

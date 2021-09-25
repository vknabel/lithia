package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseParameterDeclaration() (*ast.DeclParameter, []SyntaxError) {
	name := ast.Identifier(fp.Node.Content(fp.Source))
	return ast.MakeDeclParameter(name, fp.AstSource()), nil
}

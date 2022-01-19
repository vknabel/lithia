package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseParameterDeclaration() (*ast.DeclParameter, []SyntaxError) {
	name := ast.Identifier(fp.Node.Content(fp.Source))
	param := ast.MakeDeclParameter(name, fp.AstSource())
	param.Docs = fp.ConsumeDocs()
	return param, nil
}

func (fp *FileParser) ParseParameterDeclarationList() ([]ast.DeclParameter, []SyntaxError) {
	params := []ast.DeclParameter{}
	errors := []SyntaxError{}
	for i := 0; i < int(fp.Node.NamedChildCount()); i++ {
		child := fp.Node.NamedChild(i)
		if child.Type() == TYPE_NODE_COMMENT {
			continue
		}
		param, errs := fp.ChildParserConsumingComments(child).ParseParameterDeclaration()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		params = append(params, *param)
	}
	if len(errors) > 0 {
		return params, errors
	} else {
		return params, nil
	}
}

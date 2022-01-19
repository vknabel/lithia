package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseLetDeclaration() (*ast.DeclConstant, []SyntaxError) {
	nameIdentifier := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))

	valueNode := fp.Node.ChildByFieldName("value")
	value, errs := fp.ChildParser(valueNode).ParseExpression()
	if len(errs) != 0 {
		return nil, errs
	}
	decl := ast.MakeDeclConstant(nameIdentifier, value, fp.AstSource())
	decl.Docs = fp.ConsumeDocs()
	return decl, nil
}

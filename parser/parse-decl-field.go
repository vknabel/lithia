package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseFieldDeclaration() (*ast.DeclField, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_DATA_PROPERTY_VALUE:
		return fp.parseDataPropertyValue()
	case TYPE_NODE_DATA_PROPERTY_FUNCTION:
		return fp.parseDataPropertyFunction()
	default:
		return nil, []SyntaxError{fp.SyntaxErrorf("unexpected node type %s", fp.Node.Type())}
	}
}

func (fp *FileParser) parseDataPropertyValue() (*ast.DeclField, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	field := ast.MakeDeclField(name, nil, fp.AstSource())
	field.Docs = fp.ConsumeDocs()
	return field, nil
}

func (fp *FileParser) parseDataPropertyFunction() (*ast.DeclField, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	paramsNode := fp.Node.ChildByFieldName("parameters")

	params := make([]ast.DeclParameter, 0)
	errors := []SyntaxError{}
	for i := 0; i < int(paramsNode.NamedChildCount()); i++ {
		child := paramsNode.NamedChild(i)
		param, paramErrors := fp.SameScopeChildParser(child).ParseParameterDeclaration()
		if paramErrors != nil {
			errors = append(errors, paramErrors...)
		}
		if param != nil {
			params = append(params, *param)
		}
	}

	field := ast.MakeDeclField(name, params, fp.AstSource())
	field.Docs = fp.ConsumeDocs()
	return field, errors
}

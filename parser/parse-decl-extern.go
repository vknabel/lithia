package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseExternDeclaration() (*ast.Decl, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	parameters := fp.Node.ChildByFieldName("parameters")
	if parameters != nil {
		paramsParser := fp.ChildParser(parameters)
		var decl ast.Decl
		params, errs := paramsParser.ParseParameterDeclarationList()
		if params == nil {
			return nil, errs
		}
		funcDecl := ast.MakeDeclExternFunc(name, params, fp.AstSource())
		decl = *funcDecl
		return &decl, nil
	}

	var decl ast.Decl
	typeDecl := ast.MakeDeclExternType(name, fp.AstSource())
	decl = *typeDecl

	propertiesNode := fp.Node.ChildByFieldName("properties")
	if propertiesNode == nil {
		return &decl, nil
	}
	propsp := fp.ChildParser(propertiesNode)

	dataDecl := ast.MakeDeclData(name, fp.AstSource())
	dataDecl.Docs = fp.ConsumeDocs()

	var numberOfFields int
	if propertiesNode != nil {
		numberOfFields = int(propertiesNode.ChildCount())
	} else {
		numberOfFields = 0
	}
	errors := []SyntaxError{}
	for i := 0; i < numberOfFields; i++ {
		child := propertiesNode.NamedChild(i)
		if child.Type() == TYPE_NODE_COMMENT {
			propsp.Comments = append(propsp.Comments, child.Content(fp.Source))
		}

		childp := propsp.ChildParserConsumingComments(child)
		field, propErrors := childp.ParseFieldDeclaration()
		if propErrors != nil {
			errors = append(errors, propErrors...)
		}

		dataDecl.AddField(field)
	}
	return &decl, errors
}

package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseExternDeclaration() (*ast.Decl, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	parameters := fp.Node.ChildByFieldName("parameters")
	if parameters != nil {
		paramsParser := fp.NewScopeChildParser(parameters)
		var decl ast.Decl
		params, errs := paramsParser.ParseParameterDeclarationList()
		if params == nil {
			return nil, errs
		}
		funcDecl := ast.MakeDeclExternFunc(name, params, fp.AstSource())
		funcDecl.Docs = fp.ConsumeDocs()
		decl = *funcDecl
		return &decl, nil
	}

	typeDecl := ast.MakeDeclExternType(name, fp.AstSource())
	typeDecl.Docs = fp.ConsumeDocs()

	propertiesNode := fp.Node.ChildByFieldName("properties")
	if propertiesNode == nil {
		var decl ast.Decl
		decl = *typeDecl
		return &decl, nil
	}
	propsp := fp.NewScopeChildParser(propertiesNode)

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
			continue
		}

		childp := propsp.ChildParserConsumingComments(child)
		field, propErrors := childp.ParseFieldDeclaration()
		if len(propErrors) > 0 {
			errors = append(errors, propErrors...)
		}
		if field != nil {
			typeDecl.AddField(*field)
		}
	}
	var decl ast.Decl
	decl = *typeDecl
	return &decl, errors
}

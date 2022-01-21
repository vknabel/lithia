package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseDataDeclaration() (*ast.DeclData, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))

	dataDecl := ast.MakeDeclData(name, fp.AstSource())
	dataDecl.Docs = fp.ConsumeDocs()

	propertiesNode := fp.Node.ChildByFieldName("properties")
	if propertiesNode == nil {
		return dataDecl, nil
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
			dataDecl.AddField(*field)
		}
	}

	return dataDecl, errors
}

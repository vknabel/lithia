package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseDataDeclaration() (*ast.DeclData, []SyntaxError) {
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	propertiesNode := fp.Node.ChildByFieldName("properties")
	propsp := fp.ChildParser(propertiesNode)

	dataDecl := ast.MakeDeclData(name, fp.AstSource())
	dataDecl.Docs = fp.ConsumeDocs()

	errors := []SyntaxError{}
	for i := 0; i < int(propertiesNode.NamedChildCount()); i++ {
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

	return dataDecl, errors
}
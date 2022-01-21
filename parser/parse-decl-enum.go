package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseEnumDeclaration() (*ast.DeclEnum, []ast.Decl, []SyntaxError) {
	enumName := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	decl := ast.MakeDeclEnum(enumName, fp.AstSource())

	caseList := fp.Node.ChildByFieldName("cases")
	if caseList == nil {
		return decl, nil, nil
	}
	casep := fp.NewScopeChildParser(caseList)
	decl.Docs = fp.ConsumeDocs()

	allChildDecls := []ast.Decl{}
	errors := []SyntaxError{}

	var caseCount int
	if caseList != nil {
		caseCount = int(caseList.NamedChildCount())
	} else {
		caseCount = 0
	}
	for i := 0; i < caseCount; i++ {
		childNode := caseList.NamedChild(i)
		if childNode.Type() == TYPE_NODE_COMMENT {
			casep.Comments = append(casep.Comments, childNode.Content(fp.Source))
			continue
		}

		docs := ast.MakeDocs(fp.Comments)
		caseDecl, childDecls, err := fp.ChildParserConsumingComments(childNode).ParseEnumCaseDeclaration()
		if err != nil {
			errors = append(errors, err...)
		}
		if childDecls != nil {
			allChildDecls = append(allChildDecls, childDecls...)
		}
		if caseDecl != nil {
			caseDecl.Docs = docs
			decl.AddCase(caseDecl)
		}
	}
	return decl, allChildDecls, errors
}

package parser

import (
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseImportDeclaration() (*ast.DeclImport, []SyntaxError) {
	importModuleNode := fp.Node.ChildByFieldName("name")
	membersNode := fp.Node.ChildByFieldName("members")

	modulePath := make([]string, 0, importModuleNode.NamedChildCount())
	for i := 0; i < int(importModuleNode.NamedChildCount()); i++ {
		modulePath = append(modulePath, importModuleNode.NamedChild(i).Content(fp.Source))
	}
	moduleName := ast.ModuleName(strings.Join(modulePath, "."))
	importDecl := ast.MakeDeclImport(moduleName, fp.AstSource())

	for i := 0; i <= int(membersNode.NamedChildCount()); i++ {
		child := membersNode.NamedChild(i)
		if child.Type() == TYPE_NODE_IDENTIFIER {
			name := ast.Identifier(child.Content(fp.Source))
			member := ast.MakeDeclImportMember(name, fp.ChildParser(child).AstSource())
			importDecl.AddMember(member)
		}
	}
	return importDecl, nil
}

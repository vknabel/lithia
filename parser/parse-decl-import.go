package parser

import (
	"strings"

	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseImportDeclaration() (*ast.DeclImport, []SyntaxError) {
	importModuleNode := fp.Node.ChildByFieldName("name")
	aliasNode := fp.Node.ChildByFieldName("alias")
	membersNode := fp.Node.ChildByFieldName("members")
	var membersCount int
	if membersNode == nil {
		membersCount = 0
	} else {
		membersCount = int(membersNode.NamedChildCount())
	}

	modulePath := make([]string, 0, importModuleNode.NamedChildCount())
	for i := 0; i < int(importModuleNode.NamedChildCount()); i++ {
		modulePath = append(modulePath, importModuleNode.NamedChild(i).Content(fp.Source))
	}
	moduleName := ast.ModuleName(strings.Join(modulePath, "."))

	var importDecl *ast.DeclImport
	if aliasNode == nil {
		importDecl = ast.MakeDeclImport(moduleName, fp.AstSource())
	} else {
		importDecl = ast.MakeDeclAliasImport(ast.Identifier(aliasNode.Content(fp.Source)), moduleName, fp.AstSource())
	}

	for i := 0; i < membersCount; i++ {
		child := membersNode.NamedChild(i)
		if child.Type() == TYPE_NODE_IDENTIFIER {
			name := ast.Identifier(child.Content(fp.Source))
			member := ast.MakeDeclImportMember(moduleName, name, fp.SameScopeChildParser(child).AstSource())
			importDecl.AddMember(member)
		}
	}
	return importDecl, nil
}

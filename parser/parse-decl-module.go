package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseModuleDeclaration() (*ast.DeclModule, []SyntaxError) {
	internalName := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	decl := ast.MakeDeclModule(internalName, fp.AstSource())
	decl.Docs = fp.ConsumeDocs()
	return decl, nil
}

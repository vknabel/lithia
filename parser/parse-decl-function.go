package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseFunctionDeclaration() (*ast.DeclFunc, []SyntaxError) {
	functionNode := fp.Node.ChildByFieldName("function")
	function, errs := fp.ChildParser(functionNode).ParseFunctionExpr()

	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	funcDecl := ast.MakeDeclFunc(name, function, fp.AstSource())
	funcDecl.Docs = fp.ConsumeDocs()

	return funcDecl, errs
}

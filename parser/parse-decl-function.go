package parser

import "github.com/vknabel/lithia/ast"

func (fp *FileParser) ParseFunctionDeclaration() (*ast.DeclFunc, []SyntaxError) {
	functionNode := fp.Node.ChildByFieldName("function")
	name := ast.Identifier(fp.Node.ChildByFieldName("name").Content(fp.Source))
	function, errs := fp.NewScopeChildParser(functionNode).ParseFunctionExpr(string(name))

	funcDecl := ast.MakeDeclFunc(name, function, fp.AstSource())
	funcDecl.Docs = fp.ConsumeDocs()

	return funcDecl, errs
}

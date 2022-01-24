package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseEnumCaseDeclaration() (*ast.DeclEnumCase, []ast.Decl, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_ENUM_CASE_REFERENCE:
		return fp.parseEnumCaseReference()
	case TYPE_NODE_DATA_DECLARATION:
		dataDecl, errors := fp.ParseDataDeclaration()
		return ast.MakeDeclEnumCase(dataDecl.Name), []ast.Decl{*dataDecl}, errors
	case TYPE_NODE_ENUM_DECLARATION:
		enumDecl, childDecls, errors := fp.ParseEnumDeclaration()
		enumCase := ast.MakeDeclEnumCase(enumDecl.Name)
		return enumCase, append(childDecls, *enumDecl), errors
	default:
		return nil, nil, []SyntaxError{fp.SyntaxErrorf("unexpected node")}
	}
}

func (fp *FileParser) parseEnumCaseReference() (*ast.DeclEnumCase, []ast.Decl, []SyntaxError) {
	enumCase := ast.MakeDeclEnumCase(ast.Identifier(fp.Node.Content(fp.Source)))
	enumCase.Docs = fp.ConsumeDocs()
	return enumCase, nil, nil
}

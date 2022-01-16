package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseDeclsIfGiven() ([]ast.Decl, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_MODULE_DECLARATION:
		stmt, err := fp.ParseModuleDeclaration()
		return []ast.Decl{*stmt}, err
	case TYPE_NODE_IMPORT_DECLARATION:
		stmt, err := fp.ParseImportDeclaration()
		return []ast.Decl{*stmt}, err
	case TYPE_NODE_LET_DECLARATION:
		stmt, err := fp.ParseLetDeclaration()
		return []ast.Decl{*stmt}, err
	case TYPE_NODE_ENUM_DECLARATION:
		stmt, childDecls, err := fp.ParseEnumDeclaration()
		return append(childDecls, *stmt), err
	case TYPE_NODE_DATA_DECLARATION:
		stmt, err := fp.ParseDataDeclaration()
		return []ast.Decl{*stmt}, err
	case TYPE_NODE_FUNCTION_DECLARATION:
		stmt, err := fp.ParseFunctionDeclaration()
		return []ast.Decl{*stmt}, err
	case TYPE_NODE_EXTERN_DECLARATION:
		stmt, err := fp.ParseExternDeclaration()
		return []ast.Decl{*stmt}, err

	default:
		return nil, nil
	}
}

package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseDeclsIfGiven() ([]ast.Decl, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_MODULE_DECLARATION:
		stmt, err := fp.ParseModuleDeclaration()
		if stmt != nil {
			return []ast.Decl{*stmt}, err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_IMPORT_DECLARATION:
		importDecl, err := fp.ParseImportDeclaration()
		if importDecl != nil {
			memberImportDecls := make([]ast.Decl, len(importDecl.Members))
			for i, member := range importDecl.Members {
				memberImportDecls[i] = member
			}
			return append(memberImportDecls, *importDecl), err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_LET_DECLARATION:
		stmt, err := fp.ParseLetDeclaration()
		if stmt != nil {
			return []ast.Decl{*stmt}, err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_ENUM_DECLARATION:
		stmt, childDecls, err := fp.ParseEnumDeclaration()
		if stmt != nil {
			return append(childDecls, *stmt), err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_DATA_DECLARATION:
		stmt, err := fp.ParseDataDeclaration()
		if stmt != nil {
			return []ast.Decl{*stmt}, err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_FUNCTION_DECLARATION:
		stmt, err := fp.ParseFunctionDeclaration()
		if stmt != nil {
			return []ast.Decl{*stmt}, err
		} else {
			return []ast.Decl{}, err
		}
	case TYPE_NODE_EXTERN_DECLARATION:
		stmt, err := fp.ParseExternDeclaration()
		if stmt != nil {
			return []ast.Decl{*stmt}, err
		} else {
			return []ast.Decl{}, err
		}

	default:
		return nil, nil
	}
}

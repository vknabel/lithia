package ast

import "fmt"

var _ Decl = DeclImportMember{}
var _ Overviewable = DeclImportMember{}

type DeclImportMember struct {
	Name       Identifier
	ModuleName ModuleName

	MetaInfo *MetaDecl
}

func (e DeclImportMember) DeclName() Identifier {
	return e.Name
}

func (e DeclImportMember) DeclOverview() string {
	return fmt.Sprintf("import %s { %s }", e.ModuleName, e.Name)
}

func (e DeclImportMember) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclImportMember) IsExportedDecl() bool {
	return false
}

func MakeDeclImportMember(moduleName ModuleName, name Identifier, source *Source) DeclImportMember {
	return DeclImportMember{
		Name:       name,
		ModuleName: moduleName,
		MetaInfo:   &MetaDecl{source},
	}
}

func (DeclImportMember) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

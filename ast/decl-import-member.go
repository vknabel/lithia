package ast

var _ Decl = DeclImportMember{}

type DeclImportMember struct {
	Name       Identifier
	ModuleName ModuleName

	MetaInfo *MetaDecl
}

func (e DeclImportMember) DeclName() Identifier {
	return e.Name
}

func (e DeclImportMember) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclImportMember) IsExportedDecl() bool {
	return true
}

func MakeDeclImportMember(moduleName ModuleName, name Identifier, source *Source) DeclImportMember {
	return DeclImportMember{
		Name:       name,
		ModuleName: moduleName,
		MetaInfo:   &MetaDecl{source},
	}
}

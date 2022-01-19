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

func MakeDeclImportMember(moduleName ModuleName, name Identifier, source *Source) DeclImportMember {
	return DeclImportMember{
		Name:       name,
		ModuleName: moduleName,
		MetaInfo:   &MetaDecl{source},
	}
}

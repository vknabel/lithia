package ast

var _ Decl = DeclImportMember{}

type DeclImportMember struct {
	Name Identifier

	MetaInfo *MetaDecl
}

func (e DeclImportMember) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclImportMember(name Identifier, source *Source) *DeclImportMember {
	return &DeclImportMember{
		Name:     name,
		MetaInfo: &MetaDecl{source},
	}
}

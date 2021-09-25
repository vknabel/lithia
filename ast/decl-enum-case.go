package ast

var _ Decl = DeclEnumCase{}

type DeclEnumCase struct {
	Name Identifier

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclEnumCase) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclEnumCase(name Identifier) *DeclEnumCase {
	return &DeclEnumCase{
		Name:     name,
		MetaInfo: &MetaDecl{},
	}
}

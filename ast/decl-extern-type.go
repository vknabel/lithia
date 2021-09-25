package ast

var _ Decl = DeclExternType{}

type DeclExternType struct {
	Name   Identifier
	Fields map[Identifier]*DeclField

	MetaInfo *MetaDecl
}

func (e DeclExternType) Meta() *MetaDecl {
	return e.MetaInfo
}

package ast

var _ Decl = DeclExternFunc{}

type DeclExternFunc struct {
	Name       Identifier
	Parameters []*DeclParameter

	MetaInfo *MetaDecl
}

func (e DeclExternFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

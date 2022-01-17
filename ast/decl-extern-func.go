package ast

var _ Decl = DeclExternFunc{}

type DeclExternFunc struct {
	Name       Identifier
	Parameters []DeclParameter

	MetaInfo *MetaDecl
}

func (e DeclExternFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclExternFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclExternFunc(name Identifier, params []DeclParameter, source *Source) *DeclExternFunc {
	return &DeclExternFunc{
		Name:       name,
		Parameters: params,
		MetaInfo:   &MetaDecl{source},
	}
}

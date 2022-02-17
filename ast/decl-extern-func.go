package ast

var _ Decl = DeclExternFunc{}

type DeclExternFunc struct {
	Name       Identifier
	Parameters []DeclParameter

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclExternFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclExternFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclExternFunc) IsExportedDecl() bool {
	return true
}

func MakeDeclExternFunc(name Identifier, params []DeclParameter, source *Source) *DeclExternFunc {
	return &DeclExternFunc{
		Name:       name,
		Parameters: params,
		MetaInfo:   &MetaDecl{source},
	}
}

func (decl DeclExternFunc) ProvidedDocs() *Docs {
	return decl.Docs
}

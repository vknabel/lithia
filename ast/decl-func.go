package ast

var _ Decl = DeclFunc{}

type DeclFunc struct {
	Name Identifier
	Impl *ExprFunc

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclFunc) DeclName() Identifier {
	return e.Name
}

func (e DeclFunc) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclFunc) IsExportedDecl() bool {
	return true
}

func MakeDeclFunc(name Identifier, impl *ExprFunc, source *Source) *DeclFunc {
	return &DeclFunc{
		Name: name,
		Impl: impl,
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (decl DeclFunc) ProvidedDocs() *Docs {
	return decl.Docs
}

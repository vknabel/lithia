package ast

var _ Decl = DeclFunc{}

type DeclFunc struct {
	Name Identifier
	Impl *ExprFunc

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclFunc) Meta() *MetaDecl {
	return e.MetaInfo
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

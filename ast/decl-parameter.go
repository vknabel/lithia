package ast

var _ Decl = DeclParameter{}

type DeclParameter struct {
	Name Identifier

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclParameter) DeclName() Identifier {
	return e.Name
}

func (e DeclParameter) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclParameter(name Identifier, source *Source) *DeclParameter {
	return &DeclParameter{
		Name: name,
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

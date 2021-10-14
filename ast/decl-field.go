package ast

var _ Decl = DeclField{}

type DeclField struct {
	Name       Identifier
	Parameters []*DeclParameter

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclField) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclField(name Identifier, params []*DeclParameter, source *Source) *DeclField {
	return &DeclField{
		Name:       name,
		Parameters: params,
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}
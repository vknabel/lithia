package ast

var _ Decl = DeclData{}

type DeclData struct {
	Name   Identifier
	Fields map[Identifier]*DeclField

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclData) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclData(name Identifier, source *Source) *DeclData {
	return &DeclData{
		Name:   name,
		Fields: map[Identifier]*DeclField{},
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (e *DeclData) AddField(field *DeclField) {
	e.Fields[field.Name] = field
}

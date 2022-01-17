package ast

var _ Decl = DeclData{}

type DeclData struct {
	Name   Identifier
	Fields []*DeclField

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclData) DeclName() Identifier {
	return e.Name
}

func (e DeclData) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclData(name Identifier, source *Source) *DeclData {
	return &DeclData{
		Name:   name,
		Fields: []*DeclField{},
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (e *DeclData) AddField(field *DeclField) {
	e.Fields = append(e.Fields, field)
}

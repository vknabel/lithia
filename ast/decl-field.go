package ast

var _ Decl = DeclField{}

type DeclField struct {
	Name       Identifier
	Parameters []DeclParameter

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclField) DeclName() Identifier {
	return e.Name
}

func (e DeclField) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclField) IsExportedDecl() bool {
	return true
}

func MakeDeclField(name Identifier, params []DeclParameter, source *Source) *DeclField {
	return &DeclField{
		Name:       name,
		Parameters: params,
		Docs:       MakeDocs([]string{}),
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (decl DeclField) ProvidedDocs() *Docs {
	return decl.Docs
}

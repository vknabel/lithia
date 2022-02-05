package ast

var _ Decl = DeclExternType{}

type DeclExternType struct {
	Name   Identifier
	Fields map[Identifier]DeclField

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclExternType) DeclName() Identifier {
	return e.Name
}

func (e DeclExternType) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclExternType) IsExportedDecl() bool {
	return true
}

func MakeDeclExternType(name Identifier, source *Source) *DeclExternType {
	return &DeclExternType{
		Name:     name,
		Fields:   make(map[Identifier]DeclField),
		Docs:     nil,
		MetaInfo: &MetaDecl{source},
	}
}

func (e *DeclExternType) AddField(decl DeclField) {
	e.Fields[decl.Name] = decl
}

func (decl DeclExternType) ProvidedDocs() *Docs {
	return decl.Docs
}

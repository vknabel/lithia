package ast

var _ Decl = DeclModule{}

type DeclModule struct {
	Name Identifier

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclModule) DeclName() Identifier {
	return e.Name
}

func (e DeclModule) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclModule) IsExportedDecl() bool {
	return false
}

func MakeDeclModule(internalName Identifier, source *Source) *DeclModule {
	return &DeclModule{Name: internalName, MetaInfo: &MetaDecl{source}}
}

func (decl DeclModule) ProvidedDocs() *Docs {
	return decl.Docs
}

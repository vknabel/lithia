package ast

var _ Decl = DeclModule{}

type DeclModule struct {
	Name Identifier

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclModule) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclModule(internalName Identifier, source *Source) *DeclModule {
	return &DeclModule{Name: internalName, MetaInfo: &MetaDecl{source}}
}

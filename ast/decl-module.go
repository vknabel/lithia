package ast

import "fmt"

var _ Decl = DeclModule{}
var _ Overviewable = DeclModule{}

type DeclModule struct {
	Name Identifier

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclModule) DeclName() Identifier {
	return e.Name
}

func (e DeclModule) DeclOverview() string {
	return fmt.Sprintf("module %s", e.Name)
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

func (DeclModule) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

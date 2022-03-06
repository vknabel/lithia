package ast

type Decl interface {
	DeclName() Identifier
	IsExportedDecl() bool
	Meta() *MetaDecl

	EnumerateNestedDecls(enumerate func(interface{}, []Decl))
}

package ast

type Decl interface {
	DeclName() Identifier
	IsExportedDecl() bool
	Meta() *MetaDecl
}

package ast

type Decl interface {
	DeclName() Identifier
	Meta() *MetaDecl
}

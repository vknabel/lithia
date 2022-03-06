package ast

type Expr interface {
	Meta() *MetaExpr
	EnumerateNestedDecls(enumerate func(interface{}, []Decl))
}

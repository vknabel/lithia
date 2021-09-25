package ast

var _ Expr = ExprString{}

type ExprString struct {
	Literal string

	MetaInfo *MetaExpr
}

func (e ExprString) Meta() *MetaExpr {
	return e.MetaInfo
}

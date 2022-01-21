package ast

var _ Expr = ExprFloat{}

type ExprFloat struct {
	Literal float64

	MetaInfo *MetaExpr
}

func (e ExprFloat) Meta() *MetaExpr {
	return e.MetaInfo
}
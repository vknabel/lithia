package ast

var _ Expr = ExprGroup{}

type ExprGroup struct {
	Expr

	MetaInfo *MetaExpr
}

func (e ExprGroup) Meta() *MetaExpr {
	return e.MetaInfo
}

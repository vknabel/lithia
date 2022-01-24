package ast

var _ Expr = ExprOperatorUnary{}

type ExprOperatorUnary struct {
	Operator OperatorUnary
	Expr     *Expr

	MetaInfo *MetaExpr
}

func (e ExprOperatorUnary) Meta() *MetaExpr {
	return e.MetaInfo
}

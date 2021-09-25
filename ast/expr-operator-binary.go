package ast

var _ Expr = ExprOperatorBinary{}

type ExprOperatorBinary struct {
	Operator OperatorBinary
	Left     *Expr
	Right    *Expr

	MetaInfo *MetaExpr
}

func (e ExprOperatorBinary) Meta() *MetaExpr {
	return e.MetaInfo
}

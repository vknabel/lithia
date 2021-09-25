package ast

var _ Expr = ExprInvocation{}

type ExprInvocation struct {
	Function  *Expr
	Arguments []*Expr

	MetaInfo *MetaExpr
}

func (e ExprInvocation) Meta() *MetaExpr {
	return e.MetaInfo
}

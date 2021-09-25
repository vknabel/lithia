package ast

var _ Expr = ExprArray{}

type ExprArray struct {
	Elements []*Expr

	MetaInfo *MetaExpr
}

func (e ExprArray) Meta() *MetaExpr {
	return e.MetaInfo
}

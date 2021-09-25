package ast

var _ Expr = ExprTypeSwitch{}

type ExprTypeSwitch struct {
	Type  Expr
	Cases map[Identifier]*Expr

	MetaInfo *MetaExpr
}

func (e ExprTypeSwitch) Meta() *MetaExpr {
	return e.MetaInfo
}

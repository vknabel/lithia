package ast

var _ Expr = ExprInt{}

type ExprInt struct {
	Literal int64

	MetaInfo *MetaExpr
}

func (e ExprInt) Meta() *MetaExpr {
	return e.MetaInfo
}

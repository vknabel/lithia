package ast

var _ Expr = ExprIdentifier{}

type ExprIdentifier struct {
	Name Identifier

	MetaInfo *MetaExpr
}

func (e ExprIdentifier) Meta() *MetaExpr {
	return e.MetaInfo
}

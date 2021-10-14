package ast

var _ Expr = ExprInt{}

type ExprInt struct {
	Literal int64

	MetaInfo *MetaExpr
}

func (e ExprInt) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprInt(literal int64, source *Source) *ExprInt {
	return &ExprInt{
		Literal: literal,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

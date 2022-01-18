package ast

var _ Expr = ExprOperatorBinary{}

type ExprOperatorBinary struct {
	Operator OperatorBinary
	Left     Expr
	Right    Expr

	MetaInfo *MetaExpr
}

func (e ExprOperatorBinary) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprOperatorBinary(operator OperatorBinary, left, right Expr, source *Source) *ExprOperatorBinary {
	return &ExprOperatorBinary{
		Operator: operator,
		Left:     left,
		Right:    right,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

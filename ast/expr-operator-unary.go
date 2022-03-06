package ast

var _ Expr = ExprOperatorUnary{}

type ExprOperatorUnary struct {
	Operator OperatorUnary
	Expr     Expr

	MetaInfo *MetaExpr
}

func (e ExprOperatorUnary) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprOperatorUnary(operator OperatorUnary, expr Expr, source *Source) *ExprOperatorUnary {
	return &ExprOperatorUnary{
		Operator: operator,
		Expr:     expr,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprOperatorUnary) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Expr.EnumerateNestedDecls(enumerate)
}

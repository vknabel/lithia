package ast

var _ Expr = ExprInvocation{}

type ExprInvocation struct {
	Function  Expr
	Arguments []*Expr

	MetaInfo *MetaExpr
}

func (e ExprInvocation) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprInvocation(function Expr) *ExprInvocation {
	return &ExprInvocation{
		Function: function,
		MetaInfo: &MetaExpr{},
	}
}

func (e *ExprInvocation) AddArgument(argument Expr) {
	e.Arguments = append(e.Arguments, &argument)
}

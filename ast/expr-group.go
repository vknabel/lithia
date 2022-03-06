package ast

var _ Expr = ExprGroup{}

type ExprGroup struct {
	Expr Expr

	MetaInfo *MetaExpr
}

func (e ExprGroup) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprGroup(expr Expr, source *Source) *ExprGroup {
	return &ExprGroup{
		Expr:     expr,
		MetaInfo: &MetaExpr{Source: source},
	}
}

func (e ExprGroup) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Expr.EnumerateNestedDecls(enumerate)
}

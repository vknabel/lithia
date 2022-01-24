package ast

var _ Expr = ExprArray{}

type ExprArray struct {
	Elements []Expr

	MetaInfo *MetaExpr
}

func (e ExprArray) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprArray(elements []Expr, source *Source) *ExprArray {
	return &ExprArray{
		Elements: elements,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

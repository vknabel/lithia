package ast

var _ Expr = ExprFloat{}

type ExprFloat struct {
	Literal float64

	MetaInfo *MetaExpr
}

func (e ExprFloat) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprFloat(literal float64, source *Source) *ExprFloat {
	return &ExprFloat{
		Literal: literal,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprFloat) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

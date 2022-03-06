package ast

var _ Expr = ExprString{}

type ExprString struct {
	Literal string

	MetaInfo *MetaExpr
}

func (e ExprString) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprString(literal string, source *Source) *ExprString {
	return &ExprString{
		Literal: literal,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprString) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

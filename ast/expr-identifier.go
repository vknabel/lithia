package ast

var _ Expr = ExprIdentifier{}

type ExprIdentifier struct {
	Name Identifier

	MetaInfo *MetaExpr
}

func (e ExprIdentifier) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprIdentifier(name Identifier, source *Source) *ExprIdentifier {
	return &ExprIdentifier{
		Name: name,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprIdentifier) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

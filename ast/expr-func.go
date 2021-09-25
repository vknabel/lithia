package ast

var _ Expr = ExprFunc{}

type ExprFunc struct {
	Parameters   []DeclParameter
	Declarations []*Decl
	Statements   []*Expr

	MetaInfo *MetaExpr
}

func (e ExprFunc) Meta() *MetaExpr {
	return e.MetaInfo
}

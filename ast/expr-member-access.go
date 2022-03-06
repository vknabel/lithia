package ast

var _ Expr = ExprMemberAccess{}

type ExprMemberAccess struct {
	Target     Expr
	AccessPath []Identifier

	MetaInfo *MetaExpr
}

func (e ExprMemberAccess) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprMemberAccess(target Expr, accessPath []Identifier, source *Source) *ExprMemberAccess {
	return &ExprMemberAccess{
		Target:     target,
		AccessPath: accessPath,
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e ExprMemberAccess) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	e.Target.EnumerateNestedDecls(enumerate)
}

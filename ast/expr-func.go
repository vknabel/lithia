package ast

var _ Expr = ExprFunc{}

type ExprFunc struct {
	Name         string
	Parameters   []DeclParameter
	Declarations []Decl
	Expressions  []Expr

	MetaInfo *MetaExpr
}

func (e ExprFunc) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprFunc(name string, parameters []DeclParameter, source *Source) *ExprFunc {
	return &ExprFunc{
		Name:         name,
		Parameters:   parameters,
		Declarations: []Decl{},
		Expressions:  []Expr{},
		MetaInfo:     &MetaExpr{Source: source},
	}
}

func (e *ExprFunc) AddDecl(decl Decl) {
	e.Declarations = append(e.Declarations, decl)
}

func (e *ExprFunc) AddExpr(expr Expr) {
	e.Expressions = append(e.Expressions, expr)
}

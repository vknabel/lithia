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

func (e ExprFunc) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	enumeratedDecls := make([]Decl, len(e.Parameters)+len(e.Declarations))
	for i, param := range e.Parameters {
		enumeratedDecls[i] = param
	}
	for i, decl := range e.Declarations {
		decl.EnumerateNestedDecls(enumerate)
		enumeratedDecls[len(e.Parameters)+i] = decl
	}
	enumerate(e, enumeratedDecls)

	for _, expr := range e.Expressions {
		expr.EnumerateNestedDecls(enumerate)
	}
}

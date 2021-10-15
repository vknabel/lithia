package ast

var _ Expr = ExprTypeSwitch{}

type ExprTypeSwitch struct {
	Type  Expr
	Cases map[Identifier]*Expr

	MetaInfo *MetaExpr
}

func (e ExprTypeSwitch) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprTypeSwitch(type_ Expr, source *Source) *ExprTypeSwitch {
	return &ExprTypeSwitch{
		Type:  type_,
		Cases: make(map[Identifier]*Expr),
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e *ExprTypeSwitch) AddCase(key Identifier, value *Expr) {
	e.Cases[key] = value
}

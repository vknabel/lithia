package ast

var _ Expr = ExprTypeSwitch{}

type ExprTypeSwitch struct {
	Type      Expr
	CaseOrder []Identifier
	Cases     map[Identifier]Expr

	MetaInfo *MetaExpr
}

func (e ExprTypeSwitch) Meta() *MetaExpr {
	return e.MetaInfo
}

func MakeExprTypeSwitch(type_ Expr, source *Source) *ExprTypeSwitch {
	return &ExprTypeSwitch{
		Type:      type_,
		CaseOrder: make([]Identifier, 0),
		Cases:     make(map[Identifier]Expr),
		MetaInfo: &MetaExpr{
			Source: source,
		},
	}
}

func (e *ExprTypeSwitch) AddCase(key Identifier, value Expr) {
	e.CaseOrder = append(e.CaseOrder, key)
	e.Cases[key] = value
}

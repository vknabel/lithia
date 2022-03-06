package ast

type MetaExpr struct {
	Source *Source
}

func (m MetaExpr) SourceLocation() *Source {
	return m.Source
}

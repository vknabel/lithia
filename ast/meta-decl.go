package ast

type MetaDecl struct {
	*Source
}

func (m MetaDecl) SourceLocation() *Source {
	return m.Source
}

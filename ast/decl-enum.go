package ast

var _ Decl = DeclEnum{}

type DeclEnum struct {
	Name  Identifier
	Cases []*DeclEnumCase

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclEnum) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclEnum(name Identifier, source *Source) *DeclEnum {
	return &DeclEnum{
		Name:  name,
		Cases: []*DeclEnumCase{},
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (e *DeclEnum) AddCase(case_ *DeclEnumCase) {
	e.Cases = append(e.Cases, case_)
}

package ast

import "fmt"

var _ Decl = DeclEnum{}

type DeclEnum struct {
	Name  Identifier
	Cases []*DeclEnumCase

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclEnum) DeclName() Identifier {
	return e.Name
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

func (e DeclEnum) String() string {
	declarationClause := fmt.Sprintf("enum %s", e.Name)
	if len(e.Cases) == 0 {
		return declarationClause
	}
	declarationClause += " { "
	for _, caseDecl := range e.Cases {
		declarationClause += string(caseDecl.Name) + "; "
	}
	return declarationClause + "}"
}

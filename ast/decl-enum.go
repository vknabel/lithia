package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclEnum{}
var _ Overviewable = DeclEnum{}

type DeclEnum struct {
	Name  Identifier
	Cases []*DeclEnumCase

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclEnum) DeclName() Identifier {
	return e.Name
}

func (e DeclEnum) DeclOverview() string {
	if len(e.Cases) == 0 {
		return fmt.Sprintf("enum %s", e.Name)
	}
	caseLines := make([]string, 0)
	for _, cs := range e.Cases {
		caseLines = append(caseLines, "    "+string(cs.Name))
	}
	return fmt.Sprintf("enum %s {\n%s\n}", e.Name, strings.Join(caseLines, "\n"))
}

func (e DeclEnum) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclEnum) IsExportedDecl() bool {
	return true
}

func MakeDeclEnum(name Identifier, source *Source) *DeclEnum {
	return &DeclEnum{
		Name:  name,
		Cases: []*DeclEnumCase{},
		Docs:  MakeDocs([]string{}),
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

func (decl DeclEnum) ProvidedDocs() *Docs {
	return decl.Docs
}

func (decl DeclEnum) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls - will be handled by the parser
}

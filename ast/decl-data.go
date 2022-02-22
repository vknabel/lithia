package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclData{}
var _ Overviewable = DeclData{}

type DeclData struct {
	Name   Identifier
	Fields []DeclField

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclData) DeclName() Identifier {
	return e.Name
}

func (e DeclData) DeclOverview() string {
	if len(e.Fields) == 0 {
		return fmt.Sprintf("data %s", e.Name)
	}
	fieldLines := make([]string, 0)
	for _, field := range e.Fields {
		fieldLines = append(fieldLines, "    "+field.DeclOverview())
	}
	return fmt.Sprintf("data %s {\n%s\n}", e.Name, strings.Join(fieldLines, "\n"))
}

func (e DeclData) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclData) IsExportedDecl() bool {
	return true
}

func MakeDeclData(name Identifier, source *Source) *DeclData {
	return &DeclData{
		Name:   name,
		Fields: []DeclField{},
		Docs:   MakeDocs([]string{}),
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (e *DeclData) AddField(field DeclField) {
	e.Fields = append(e.Fields, field)
}

func (decl DeclData) ProvidedDocs() *Docs {
	return decl.Docs
}

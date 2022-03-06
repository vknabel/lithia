package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclExternType{}
var _ Overviewable = DeclExternType{}

type DeclExternType struct {
	Name   Identifier
	Fields map[Identifier]DeclField

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclExternType) DeclName() Identifier {
	return e.Name
}

func (e DeclExternType) DeclOverview() string {
	if len(e.Fields) == 0 {
		return fmt.Sprintf("extern %s", e.Name)
	}
	fieldLines := make([]string, 0)
	for _, field := range e.Fields {
		fieldLines = append(fieldLines, "    "+field.DeclOverview())
	}
	return fmt.Sprintf("extern %s {\n%s\n}", e.Name, strings.Join(fieldLines, "\n"))
}

func (e DeclExternType) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclExternType) IsExportedDecl() bool {
	return true
}

func MakeDeclExternType(name Identifier, source *Source) *DeclExternType {
	return &DeclExternType{
		Name:     name,
		Fields:   make(map[Identifier]DeclField),
		Docs:     nil,
		MetaInfo: &MetaDecl{source},
	}
}

func (e *DeclExternType) AddField(decl DeclField) {
	e.Fields[decl.Name] = decl
}

func (decl DeclExternType) ProvidedDocs() *Docs {
	return decl.Docs
}

func (DeclExternType) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

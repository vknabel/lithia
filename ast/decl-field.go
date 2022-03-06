package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclField{}
var _ Overviewable = DeclField{}

type DeclField struct {
	Name       Identifier
	Parameters []DeclParameter

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclField) DeclName() Identifier {
	return e.Name
}

func (e DeclField) DeclOverview() string {
	if len(e.Parameters) == 0 {
		return string(e.Name)
	}
	paramNames := make([]string, len(e.Parameters))
	for i, param := range e.Parameters {
		paramNames[i] = string(param.Name)
	}
	return fmt.Sprintf("%s %s", e.Name, strings.Join(paramNames, ", "))
}

func (e DeclField) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclField) IsExportedDecl() bool {
	return true
}

func MakeDeclField(name Identifier, params []DeclParameter, source *Source) *DeclField {
	return &DeclField{
		Name:       name,
		Parameters: params,
		Docs:       MakeDocs([]string{}),
		MetaInfo: &MetaDecl{
			Source: source,
		},
	}
}

func (decl DeclField) ProvidedDocs() *Docs {
	return decl.Docs
}

func (DeclField) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

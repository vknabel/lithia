package ast

import (
	"fmt"
	"strings"
)

var _ Decl = DeclImport{}
var _ Overviewable = DeclImport{}

type DeclImport struct {
	Alias      Identifier
	ModuleName ModuleName
	Members    []DeclImportMember

	MetaInfo *MetaDecl
}

func (e DeclImport) DeclName() Identifier {
	return e.Alias
}

func (e DeclImport) DeclOverview() string {
	return fmt.Sprintf("import %s", e.ModuleName)
}

func (e DeclImport) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclImport) IsExportedDecl() bool {
	return false
}

func (e *DeclImport) AddMember(member DeclImportMember) {
	e.Members = append(e.Members, member)
}

func MakeDeclImport(name ModuleName, source *Source) *DeclImport {
	segments := strings.Split(string(name), ".")
	alias := Identifier(segments[len(segments)-1])
	return &DeclImport{
		Alias:      alias,
		ModuleName: name,
		Members:    make([]DeclImportMember, 0),
		MetaInfo:   &MetaDecl{source},
	}
}

func MakeDeclAliasImport(alias Identifier, name ModuleName, source *Source) *DeclImport {
	return &DeclImport{
		Alias:      alias,
		ModuleName: name,
		Members:    make([]DeclImportMember, 0),
		MetaInfo:   &MetaDecl{source},
	}
}

func (DeclImport) EnumerateNestedDecls(enumerate func(interface{}, []Decl)) {
	// no nested decls
}

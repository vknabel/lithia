package ast

import "strings"

var _ Decl = DeclImport{}

type DeclImport struct {
	ModuleName ModuleName
	Members    []DeclImportMember

	MetaInfo *MetaDecl
}

func (e DeclImport) DeclName() Identifier {
	segments := strings.Split(string(e.ModuleName), ".")
	return Identifier(segments[len(segments)-1])
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
	return &DeclImport{
		ModuleName: name,
		Members:    make([]DeclImportMember, 0),
		MetaInfo:   &MetaDecl{source},
	}
}

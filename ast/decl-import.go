package ast

var _ Decl = DeclImport{}

type DeclImport struct {
	ModuleName ModuleName
	Members    []*DeclImportMember

	MetaInfo *MetaDecl
}

func (e DeclImport) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e *DeclImport) AddMember(member *DeclImportMember) {
	e.Members = append(e.Members, member)
}

func MakeDeclImport(name ModuleName, source *Source) *DeclImport {
	return &DeclImport{
		ModuleName: name,
		Members:    make([]*DeclImportMember, 0),
		MetaInfo:   &MetaDecl{source},
	}
}

package ast

var _ Decl = DeclExternType{}

type DeclExternType struct {
	Name   Identifier
	Fields map[Identifier]*DeclField

	MetaInfo *MetaDecl
}

func (e DeclExternType) DeclName() Identifier {
	return e.Name
}

func (e DeclExternType) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclExternType(name Identifier, source *Source) *DeclExternType {
	return &DeclExternType{
		Name:     name,
		Fields:   make(map[Identifier]*DeclField),
		MetaInfo: &MetaDecl{source},
	}
}

func (e *DeclExternType) AddField(name Identifier, decl *DeclField) {
	e.Fields[name] = decl
}

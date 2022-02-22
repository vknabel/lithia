package ast

import "fmt"

var _ Decl = DeclConstant{}
var _ Overviewable = DeclConstant{}

type DeclConstant struct {
	Name  Identifier
	Value Expr

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclConstant) DeclName() Identifier {
	return e.Name
}

func (e DeclConstant) DeclOverview() string {
	return fmt.Sprintf("let %s", e.Name)
}

func (e DeclConstant) Meta() *MetaDecl {
	return e.MetaInfo
}

func (e DeclConstant) IsExportedDecl() bool {
	return true
}

func MakeDeclConstant(name Identifier, value Expr, source *Source) *DeclConstant {
	return &DeclConstant{
		Name:     name,
		Value:    value,
		MetaInfo: &MetaDecl{source},
	}
}

func (e DeclConstant) ProvidedDocs() *Docs {
	return e.Docs
}

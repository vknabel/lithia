package ast

var _ Decl = DeclConstant{}

type DeclConstant struct {
	Name  Identifier
	Value Expr

	Docs     *Docs
	MetaInfo *MetaDecl
}

func (e DeclConstant) DeclName() Identifier {
	return e.Name
}

func (e DeclConstant) Meta() *MetaDecl {
	return e.MetaInfo
}

func MakeDeclConstant(name Identifier, value Expr, source *Source) *DeclConstant {
	return &DeclConstant{
		Name:     name,
		Value:    value,
		MetaInfo: &MetaDecl{source},
	}
}

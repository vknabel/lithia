package ast

type ModuleName string

type ContextModule struct {
	Name         ModuleName
	Imports      []ModuleName
	Declarations []*Decl
	Statements   []*Expr
}

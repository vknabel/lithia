package parser

type Start struct {
	Package string
	Imports []Import
	Stmts   []Stmt
}

type Import struct {
	Identifier string
	Members    []string
}

type Stmt interface {
}

type DataDecl struct {
	Name       string
	Properties []PropertyDecl
}

type PropertyDecl struct {
	Name       string
	Parameters []string
}

type EnumDecl struct {
	Name  string
	Cases []string
}

type NamedFuncDecl struct {
	Name     string
	Function FuncLiteral
}

type Expr struct {
	Head SimpleExpr
	Tail []SimpleExpr
}

type SimpleExpr struct {
	Literal Literal
	Members []string
}

type Grouping struct {
	Literal Literal
	Members []string
}

type Literal interface{}

type FuncLiteral struct {
	Parameters []string
	Stmts      []Stmt
}

type ArrayLiteral struct {
	Values []SimpleExpr
}

type GroupingLiteral struct {
	Expr Expr
}

type StringLiteral struct {
	Value string
}

type IntLiteral struct {
	Value int
}

type FloatLiteral struct {
	Value float64
}

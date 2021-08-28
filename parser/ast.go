package parser

type AstNode interface {
	Accept(visitor AstNodeVisitor) error
}
type AstNodeVisitor interface {
	VisitDocument(Document) error
	VisitImport(Import) error
	VisitDataDecl(DataDecl) error
	VisitPropertyDecl(PropertyDecl) error
	VisitEnumDecl(EnumDecl) error
	VisitNamedFuncDecl(NamedFuncDecl) error
	VisitExpr(Expr) error
	VisitSimpleExpr(SimpleExpr) error
	VisitMemberAccess(MemberAccess) error
	VisitFuncLiteral(FuncLiteral) error
	VisitArrayLiteral(ArrayLiteral) error
	VisitGroupingLiteral(GroupingLiteral) error
	VisitStringLiteral(StringLiteral) error
	VisitIntLiteral(IntLiteral) error
	VisitFloatLiteral(FloatLiteral) error
}

type Document struct {
	Package string
	Imports []Import
	Stmts   []Stmt
}

func (d Document) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitDocument(d)
}

type Import struct {
	Identifier string
	Members    []string
}

func (i Import) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitImport(i)
}

type Stmt interface {
}

type DataDecl struct {
	Name       string
	Properties []PropertyDecl
}

func (d DataDecl) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitDataDecl(d)
}

type PropertyDecl struct {
	Name       string
	Parameters []string
}

func (p PropertyDecl) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitPropertyDecl(p)
}

type EnumDecl struct {
	Name  string
	Cases []string
}

func (e EnumDecl) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitEnumDecl(e)
}

type NamedFuncDecl struct {
	Name     string
	Function FuncLiteral
}

func (n NamedFuncDecl) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitNamedFuncDecl(n)
}

type Expr struct {
	Head SimpleExpr
	Tail []SimpleExpr
}

func (e Expr) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitExpr(e)
}

type SimpleExpr struct {
	Literal Literal
	Members []string
}

func (s SimpleExpr) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitSimpleExpr(s)
}

type MemberAccess struct {
	Literal Literal
	Members []string
}

func (m MemberAccess) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitMemberAccess(m)
}

type Literal interface{}

type FuncLiteral struct {
	Parameters []string
	Stmts      []Stmt
}

func (f FuncLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitFuncLiteral(f)
}

type ArrayLiteral struct {
	Values []SimpleExpr
}

func (a ArrayLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitArrayLiteral(a)
}

type GroupingLiteral struct {
	Expr Expr
}

func (g GroupingLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitGroupingLiteral(g)
}

type StringLiteral struct {
	Value string
}

func (s StringLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitStringLiteral(s)
}

type IntLiteral struct {
	Value int
}

func (i IntLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitIntLiteral(i)
}

type FloatLiteral struct {
	Value float64
}

func (f FloatLiteral) Accept(visitor AstNodeVisitor) error {
	return visitor.VisitFloatLiteral(f)
}

package parser

type AstPrinter struct{}

func (AstPrinter) VisitDocument(document Document) error {
	print(document)
	return nil
}

func (AstPrinter) VisitImport(Import) error {
	return nil
}

func (AstPrinter) VisitDataDecl(DataDecl) error {
	return nil
}

func (AstPrinter) VisitPropertyDecl(PropertyDecl) error {
	return nil
}

func (AstPrinter) VisitEnumDecl(EnumDecl) error {
	return nil
}

func (AstPrinter) VisitNamedFuncDecl(NamedFuncDecl) error {
	return nil
}

func (AstPrinter) VisitExpr(Expr) error {
	return nil
}

func (AstPrinter) VisitSimpleExpr(SimpleExpr) error {
	return nil
}

func (AstPrinter) VisitMemberAccess(MemberAccess) error {
	return nil
}

func (AstPrinter) VisitFuncLiteral(FuncLiteral) error {
	return nil
}

func (AstPrinter) VisitArrayLiteral(ArrayLiteral) error {
	return nil
}

func (AstPrinter) VisitGroupingLiteral(GroupingLiteral) error {
	return nil
}

func (AstPrinter) VisitStringLiteral(StringLiteral) error {
	return nil
}

func (AstPrinter) VisitIntLiteral(IntLiteral) error {
	return nil
}

func (AstPrinter) VisitFloatLiteral(FloatLiteral) error {
	return nil
}

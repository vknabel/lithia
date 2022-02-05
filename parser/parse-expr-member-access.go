package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseExprMemberAccess() (*ast.ExprMemberAccess, []SyntaxError) {
	if fp.Node.NamedChildCount() < 2 {
		return nil, []SyntaxError{fp.SyntaxErrorf("expected at least 2 children, got %d", fp.Node.NamedChildCount())}
	}
	targetExpr, errs := fp.SameScopeChildParser(fp.Node.NamedChild(0)).ParseExpression()
	if len(errs) > 0 {
		return nil, errs
	}

	keyPath := make([]ast.Identifier, 0, fp.Node.NamedChildCount()-1)
	for i := 1; i < int(fp.Node.NamedChildCount()); i++ {
		child := fp.Node.NamedChild(i)
		keyPath = append(keyPath, ast.Identifier(child.Content(fp.Source)))
	}
	return ast.MakeExprMemberAccess(targetExpr, keyPath, fp.AstSource()), nil
}

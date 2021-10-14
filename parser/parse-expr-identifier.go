package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) parseExprIdentifier() (*ast.ExprIdentifier, []SyntaxError) {
	content := fp.Node.Content(fp.Source)
	return ast.MakeExprIdentifier(ast.Identifier(content), fp.AstSource()), nil
}

package parser

import (
	"strconv"

	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseIntExpr() (*ast.ExprInt, []SyntaxError) {
	literal := fp.Node.Content(fp.Source)
	integer, err := strconv.ParseInt(literal, 10, 64)
	if err != nil {
		return nil, []SyntaxError{*fp.SyntaxErrorOrConvert(err)}
	}
	return ast.MakeExprInt(integer, fp.AstSource()), nil
}

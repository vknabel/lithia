package parser

import (
	"strconv"

	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseFloatExpr() (*ast.ExprFloat, []SyntaxError) {
	literal := fp.Node.Content(fp.Source)
	integer, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		return nil, []SyntaxError{*fp.SyntaxErrorOrConvert(err)}
	}
	return ast.MakeExprFloat(integer, fp.AstSource()), nil
}

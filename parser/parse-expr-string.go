package parser

import (
	"strconv"

	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseExprString() (*ast.ExprString, []SyntaxError) {
	string, err := strconv.Unquote(fp.Node.Content(fp.Source))
	if err != nil {
		return nil, []SyntaxError{*fp.SyntaxErrorOrConvert(err)}
	}
	return ast.MakeExprString(string, fp.AstSource()), nil
}

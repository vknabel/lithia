package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseUnaryExpr() (*ast.ExprInt, []SyntaxError) {
	return nil, []SyntaxError{fp.SyntaxErrorf("unimplemented")}
}

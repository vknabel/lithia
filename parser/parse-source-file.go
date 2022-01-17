package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseSourceFile() (*ast.SourceFile, []SyntaxError) {
	err := fp.AssertNodeType(TYPE_NODE_SOURCE_FILE)
	if err != nil {
		return nil, err
	}

	sourceFile := ast.MakeSourceFile(fp.File, fp.AstSource())
	parsingErrors := []SyntaxError{}

	for i := 0; i < int(fp.Node.NamedChildCount()); i++ {
		child := fp.Node.NamedChild(i)
		if child.Type() == TYPE_NODE_COMMENT {
			fp.Comments = append(fp.Comments, child.Content(fp.Source))
			continue
		}
		childFp := fp.ChildParserConsumingComments(child)

		parsedDecls, errs := childFp.ParseDeclsIfGiven()
		if len(errs) > 0 {
			parsingErrors = append(parsingErrors, errs...)
			continue
		}
		if parsedDecls != nil {
			for _, decl := range parsedDecls {
				sourceFile.AddDecl(decl)
			}
			continue
		}
		expr, errs := childFp.ParseExpressionIfGiven()
		if len(errs) > 0 {
			parsingErrors = append(parsingErrors, errs...)
			continue
		}
		if expr != nil {
			sourceFile.AddExpr(expr)
			continue
		}
		parsingErrors = append(parsingErrors, fp.SyntaxErrorf("unexpected %q, expected module, import, enum, data, func, let or an expression", child.Type()))
	}
	return sourceFile, parsingErrors
}

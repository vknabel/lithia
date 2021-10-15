package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseExprTypeSwitch() (*ast.ExprTypeSwitch, []SyntaxError) {
	typeNode := fp.Node.ChildByFieldName("type")
	bodyNode := fp.Node.ChildByFieldName("body")

	if typeNode == nil || bodyNode == nil {
		return nil, []SyntaxError{fp.SyntaxErrorf("expected type and body")}
	}

	typeExpr, errs := fp.ChildParser(typeNode).ParseExpression()
	if len(errs) > 0 {
		return nil, errs
	}

	typeSwitchExpr := ast.MakeExprTypeSwitch(*typeExpr, fp.AstSource())

	caseCount := int(bodyNode.NamedChildCount())
	errs = []SyntaxError{}
	for i := 0; i < caseCount; i++ {
		typeCaseNode := bodyNode.NamedChild(i)

		if typeCaseNode.Type() == TYPE_NODE_COMMENT {
			continue
		}
		labelNode := typeCaseNode.ChildByFieldName("label")
		bodyNode := typeCaseNode.ChildByFieldName("body")
		if labelNode == nil || bodyNode == nil {
			errs = append(errs, fp.SyntaxErrorf("expected label and body"))
			continue
		}
		bodyExpr, bodyErrs := fp.ChildParser(bodyNode).ParseExpression()
		if len(bodyErrs) > 0 {
			errs = append(errs, bodyErrs...)
		}
		labelIdentifier := ast.Identifier(labelNode.Content(fp.Source))
		typeSwitchExpr.AddCase(labelIdentifier, bodyExpr)
	}
	if len(errs) > 0 {
		return typeSwitchExpr, errs
	} else {
		return typeSwitchExpr, nil
	}
}

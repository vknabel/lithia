package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseExprDict() (*ast.ExprDict, []SyntaxError) {
	numberOfEntries := int(fp.Node.NamedChildCount())
	entries := make([]ast.ExprDictEntry, 0, numberOfEntries)
	for i := 0; i < numberOfEntries; i++ {
		entryNode := fp.Node.NamedChild(i)
		if entryNode.Type() == TYPE_NODE_COMMENT {
			fp.Comments = append(fp.Comments, entryNode.Content(fp.Source))
			continue
		}
		entry, errs := fp.ChildParserConsumingComments(entryNode).parseExprDictEntry()
		if len(errs) > 0 {
			return nil, errs
		}
		if entry != nil {
			entries = append(entries, *entry)
		}
	}
	return ast.MakeExprDict(entries, fp.AstSource()), nil
}

func (fp *FileParser) parseExprDictEntry() (*ast.ExprDictEntry, []SyntaxError) {
	keyNode := fp.Node.ChildByFieldName("key")
	valueNode := fp.Node.ChildByFieldName("value")

	keyP := fp.ChildParserConsumingComments(keyNode)
	key, errs := keyP.ParseExpression()
	if len(errs) > 0 {
		return nil, errs
	}

	valueP := fp.ChildParserConsumingComments(valueNode)
	value, errs := valueP.ParseExpression()
	if len(errs) > 0 {
		return nil, errs
	}
	return ast.MakeExprDictEntry(key, value, fp.AstSource()), nil
}

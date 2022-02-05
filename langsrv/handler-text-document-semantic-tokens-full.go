package langsrv

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	rc := NewReqContext(params.TextDocument)
	fileParser, err := rc.createFileParser()
	if err != nil && fileParser == nil {
		return nil, err
	}
	rootNode := fileParser.Tree.RootNode()
	tokens := highlightedTokensEntriesForNode(rootNode)
	return &protocol.SemanticTokens{
		Data: serializeHighlightedTokens(tokens),
	}, nil
}

func highlightedTokensEntriesForNode(node *sitter.Node) []highlightedToken {
	tokens := make([]highlightedToken, 0)
	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		switch child.Type() {
		case parser.TYPE_NODE_MODULE_DECLARATION:
			nameChild := child.ChildByFieldName("name")
			if nameChild != nil {
				tokens = append(tokens, highlightedToken{
					line:           uint32(nameChild.StartPoint().Row),
					column:         uint32(nameChild.StartPoint().Column),
					length:         nameChild.EndByte() - nameChild.StartByte(),
					tokenType:      token_namespace,
					tokenModifiers: []tokenModifier{modifier_declaration},
				})
			}
			keywordChild := child.Child(0)
			if keywordChild != nil {
				tokens = append(tokens,
					highlightedToken{
						line:           uint32(keywordChild.StartPoint().Row),
						column:         uint32(keywordChild.StartPoint().Column),
						length:         keywordChild.EndByte() - keywordChild.StartByte(),
						tokenType:      token_keyword,
						tokenModifiers: nil,
					},
				)
			}
		case parser.TYPE_NODE_NUMBER_LITERAL:
			tokens = append(tokens,
				highlightedToken{
					line:           uint32(child.StartPoint().Row),
					column:         uint32(child.StartPoint().Column),
					length:         child.EndByte() - child.StartByte(),
					tokenType:      token_number,
					tokenModifiers: nil,
				},
			)
		case parser.TYPE_NODE_STRING_LITERAL:
			tokens = append(tokens, highlightedToken{
				line:           uint32(child.StartPoint().Row),
				column:         uint32(child.StartPoint().Column),
				length:         child.EndByte() - child.StartByte(),
				tokenType:      token_string,
				tokenModifiers: nil,
			})
		case parser.TYPE_NODE_COMMENT:
			tokens = append(tokens,
				highlightedToken{
					line:           uint32(child.StartPoint().Row),
					column:         uint32(child.StartPoint().Column),
					length:         child.EndByte() - child.StartByte(),
					tokenType:      token_comment,
					tokenModifiers: nil,
				},
			)
		case parser.TYPE_NODE_DATA_DECLARATION:
			keywordChild := child.Child(0)
			if keywordChild != nil {
				tokens = append(tokens,
					highlightedToken{
						line:           uint32(keywordChild.StartPoint().Row),
						column:         uint32(keywordChild.StartPoint().Column),
						length:         keywordChild.EndByte() - keywordChild.StartByte(),
						tokenType:      token_keyword,
						tokenModifiers: nil,
					},
				)
			}
			nameChild := child.ChildByFieldName("name")
			if nameChild != nil {
				tokens = append(tokens, highlightedToken{
					line:           uint32(nameChild.StartPoint().Row),
					column:         uint32(nameChild.StartPoint().Column),
					length:         nameChild.EndByte() - nameChild.StartByte(),
					tokenType:      token_struct,
					tokenModifiers: []tokenModifier{modifier_declaration},
				})
			}
			tokens = append(tokens, highlightedTokensEntriesForNode(child)...)
		case parser.TYPE_NODE_FUNCTION_DECLARATION:
			keywordChild := child.Child(0)
			if keywordChild != nil {
				tokens = append(tokens,
					highlightedToken{
						line:           uint32(keywordChild.StartPoint().Row),
						column:         uint32(keywordChild.StartPoint().Column),
						length:         keywordChild.EndByte() - keywordChild.StartByte(),
						tokenType:      token_keyword,
						tokenModifiers: nil,
					},
				)
			}
			nameChild := child.ChildByFieldName("name")
			if nameChild != nil {
				tokens = append(tokens, highlightedToken{
					line:           uint32(nameChild.StartPoint().Row),
					column:         uint32(nameChild.StartPoint().Column),
					length:         nameChild.EndByte() - nameChild.StartByte(),
					tokenType:      token_function,
					tokenModifiers: []tokenModifier{modifier_declaration},
				})
			}
			tokens = append(tokens, highlightedTokensEntriesForNode(child)...)
		case parser.TYPE_NODE_TYPE_EXPRESSION:
			keywordChild := child.Child(0)
			if keywordChild != nil {
				tokens = append(tokens,
					highlightedToken{
						line:           uint32(keywordChild.StartPoint().Row),
						column:         uint32(keywordChild.StartPoint().Column),
						length:         keywordChild.EndByte() - keywordChild.StartByte(),
						tokenType:      token_keyword,
						tokenModifiers: nil,
					},
				)
			}
			tokens = append(tokens, highlightedTokensEntriesForNode(child)...)
		default:
			tokens = append(tokens, highlightedTokensEntriesForNode(child)...)
		}
	}
	return tokens
}

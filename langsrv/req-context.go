package langsrv

import (
	sitter "github.com/smacker/go-tree-sitter"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type ReqContext struct {
	textDocument protocol.TextDocumentIdentifier
	position     protocol.Position
	textDocumentEntry
}

func NewReqContext(textDocument protocol.TextDocumentIdentifier) *ReqContext {
	return &ReqContext{
		textDocument:      textDocument,
		textDocumentEntry: *langserver.documentCache.documents[textDocument.URI],
	}
}

func NewReqContextAtPosition(position *protocol.TextDocumentPositionParams) *ReqContext {
	return &ReqContext{
		textDocument:      position.TextDocument,
		position:          position.Position,
		textDocumentEntry: *langserver.documentCache.documents[position.TextDocument.URI],
	}
}

func (rc *ReqContext) findToken() (string, *protocol.Range, error) {
	node, err := rc.findNode()
	if err != nil || node == nil {
		return "", nil, err
	}
	contents := rc.textDocumentEntry.item.Text
	name := node.Content([]byte(contents))
	return name, &protocol.Range{
		Start: protocol.Position{
			Line:      uint32(node.StartPoint().Row),
			Character: uint32(node.StartPoint().Column),
		},
		End: protocol.Position{
			Line:      uint32(node.EndPoint().Row),
			Character: uint32(node.EndPoint().Column),
		},
	}, nil
}

func (rc *ReqContext) findNode() (*sitter.Node, error) {
	node := NodeAtPosition(rc.fileParser.Tree.RootNode(), rc.position)
	return node, nil
}

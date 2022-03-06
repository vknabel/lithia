package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)

	name, tokenRange, err := rc.findToken()
	if err != nil && tokenRange == nil {
		return nil, nil
	}

	for _, imported := range rc.accessibleDeclarations(context) {
		decl := imported.decl
		if string(decl.DeclName()) != name {
			continue
		}
		return &protocol.Hover{
			Contents: documentationMarkupContentForDecl(decl),
			Range:    tokenRange,
		}, nil
	}
	return nil, nil
}

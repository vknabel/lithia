package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)

	token, _, err := rc.findToken()
	if err != nil && token == "" {
		return nil, nil
	}

	for _, imported := range rc.accessibleDeclarations(context) {
		decl := imported.decl
		if string(decl.DeclName()) != token || decl.Meta().Source == nil {
			continue
		}
		return &[]protocol.LocationLink{
			{
				TargetURI: protocol.DocumentUri(decl.Meta().Source.FileName),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(decl.Meta().Source.Start.Line),
						Character: uint32(decl.Meta().Source.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(decl.Meta().Source.End.Line),
						Character: uint32(decl.Meta().Source.End.Column),
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(decl.Meta().Source.Start.Line),
						Character: uint32(decl.Meta().Source.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(decl.Meta().Source.End.Line),
						Character: uint32(decl.Meta().Source.End.Column),
					},
				},
			},
		}, nil
	}
	return nil, nil
}

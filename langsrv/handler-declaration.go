package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDeclaration(context *glsp.Context, params *protocol.DeclarationParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)

	token, _, err := rc.findToken()
	if err != nil {
		return nil, nil
	}

	for _, imported := range rc.accessibleDeclarations(context) {
		if string(imported.decl.DeclName()) != token || imported.decl.Meta().Source == nil {
			continue
		}
		return &[]protocol.LocationLink{
			{
				TargetURI: protocol.DocumentUri(imported.decl.Meta().Source.FileName),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(imported.decl.Meta().Source.Start.Line),
						Character: uint32(imported.decl.Meta().Source.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(imported.decl.Meta().Source.End.Line),
						Character: uint32(imported.decl.Meta().Source.End.Column),
					},
				},
			},
		}, nil
	}
	return nil, nil
}

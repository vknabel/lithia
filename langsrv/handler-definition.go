package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile, err := rc.parseSourceFile()
	if err != nil {
		return nil, err
	}
	token, _, err := rc.findToken()
	if err != nil {
		return nil, err
	}
	for _, decl := range sourceFile.Declarations {
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
						Character: uint32(decl.Meta().Source.End.Line),
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(decl.Meta().Source.Start.Line),
						Character: uint32(decl.Meta().Source.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(decl.Meta().Source.End.Line),
						Character: uint32(decl.Meta().Source.End.Line),
					},
				},
			},
		}, nil
	}
	return nil, nil
}

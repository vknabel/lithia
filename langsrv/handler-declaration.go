package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDeclaration(context *glsp.Context, params *protocol.DeclarationParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile, err := rc.parseSourceFile()
	if err != nil && sourceFile == nil {
		return nil, nil
	}
	token, _, err := rc.findToken()
	if err != nil {
		return nil, nil
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
			},
		}, nil
	}
	return nil, nil
}

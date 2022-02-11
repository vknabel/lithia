package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentTypeDefinition(context *glsp.Context, params *protocol.TypeDefinitionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile, err := rc.parseSourceFile()
	if err != nil && sourceFile == nil {
		return nil, nil
	}
	return nil, nil
}

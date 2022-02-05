package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentTypeDefinition(context *glsp.Context, params *protocol.TypeDefinitionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	_, err := rc.parseSourceFile()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

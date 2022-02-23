package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile := rc.sourceFile
	if sourceFile == nil {
		return nil, nil
	}
	name, tokenRange, err := rc.findToken()
	if err != nil && tokenRange == nil {
		return nil, nil
	}

	globalDeclarations := sourceFile.Declarations
	for _, sameModuleFile := range rc.textDocumentEntry.module.Files {
		fileUrl := "file://" + sameModuleFile
		if rc.item.URI == fileUrl {
			continue
		}
		docEntry := langserver.documentCache.documents[fileUrl]
		if docEntry == nil || docEntry.sourceFile == nil {
			continue
		}

		globalDeclarations = append(globalDeclarations, docEntry.sourceFile.ExportedDeclarations()...)
	}

	for _, decl := range globalDeclarations {
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

package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
)

func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile, err := rc.parseSourceFile()
	if err != nil {
		return nil, err
	}
	name, tokenRange, err := rc.findToken()
	if err != nil {
		return nil, err
	}
	for _, decl := range sourceFile.Declarations {
		if string(decl.DeclName()) != name {
			continue
		}
		var docs string
		if documented, ok := decl.(ast.Documented); ok {
			docs = documented.ProvidedDocs().Content + "\n"
		}
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind: protocol.MarkupKindMarkdown,
				Value: docs + "```lithia\n" + string(decl.DeclName()) + "\n```\n" +
					string(decl.Meta().ModuleName),
			},
			Range: tokenRange,
		}, nil
	}
	return nil, nil
}

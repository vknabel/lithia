package langsrv

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
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
	for _, decl := range sourceFile.Declarations {
		if string(decl.DeclName()) != name {
			continue
		}
		var docs string
		if documented, ok := decl.(ast.Documented); ok {
			docs = documented.ProvidedDocs().Content + "\n\n"
		}
		var overview string
		if overviewable, ok := decl.(ast.Overviewable); ok {
			overview = "```lithia\n" + overviewable.DeclOverview() + "\n```\n\n"
		}
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: overview + fmt.Sprintf("_module %s_\n\n", decl.Meta().ModuleName) + docs,
			},
			Range: tokenRange,
		}, nil
	}
	return nil, nil
}

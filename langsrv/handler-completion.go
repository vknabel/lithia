package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
)

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile, err := rc.parseSourceFile()
	if err != nil {
		return nil, err
	}
	completionItems := []protocol.CompletionItem{}
	for _, decl := range sourceFile.Declarations {
		var docs string
		if documented, ok := decl.(ast.Documented); ok {
			docs = documented.ProvidedDocs().Content + "\n"
		}
		var detail string
		if decl.Meta().ModuleName != "" {
			detail = string(decl.Meta().ModuleName) +
				"." +
				string(decl.DeclName())
		}
		kind := protocol.CompletionItemKindEnum
		insertText := string(decl.DeclName())
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      string(decl.DeclName()),
			Kind:       &kind,
			InsertText: &insertText,
			Detail:     &detail,
			Documentation: &protocol.MarkupContent{
				Kind: protocol.MarkupKindMarkdown,
				Value: docs + "```lithia\n" + string(decl.DeclName()) + "\n```\n" +
					string(decl.Meta().ModuleName),
			},
		})
	}
	return &completionItems, nil
}

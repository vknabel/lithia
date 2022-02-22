package langsrv

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
)

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)
	sourceFile := rc.textDocumentEntry.sourceFile
	if sourceFile == nil {
		return nil, nil
	}

	completionItems := []protocol.CompletionItem{}
	for _, decl := range sourceFile.Declarations {
		var docs string
		if documented, ok := decl.(ast.Documented); ok {
			docs = documented.ProvidedDocs().Content + "\n\n"
		}
		var overview string
		if overviewable, ok := decl.(ast.Overviewable); ok {
			overview = "```lithia\n" + overviewable.DeclOverview() + "\n```\n\n"
		}
		var detail string
		if decl.Meta().ModuleName != "" {
			detail = string(decl.Meta().ModuleName) +
				"." +
				string(decl.DeclName())
		}
		var kind protocol.CompletionItemKind
		switch decl.(type) {
		case ast.DeclEnum:
			kind = protocol.CompletionItemKindEnum
		case ast.DeclEnumCase:
			kind = protocol.CompletionItemKindEnumMember
		case ast.DeclConstant:
			kind = protocol.CompletionItemKindConstant
		case ast.DeclData:
			kind = protocol.CompletionItemKindStruct
		case ast.DeclExternFunc:
			kind = protocol.CompletionItemKindFunction
		case ast.DeclExternType:
			kind = protocol.CompletionItemKindClass
		case ast.DeclField:
			kind = protocol.CompletionItemKindField
		case ast.DeclFunc:
			kind = protocol.CompletionItemKindFunction
		case ast.DeclImport:
			kind = protocol.CompletionItemKindModule
		case ast.DeclImportMember:
			kind = protocol.CompletionItemKindValue
		case ast.DeclModule:
			kind = protocol.CompletionItemKindModule
		case ast.DeclParameter:
			kind = protocol.CompletionItemKindVariable
		default:
			kind = protocol.CompletionItemKindText
		}
		insertText := string(decl.DeclName())
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      string(decl.DeclName()),
			Kind:       &kind,
			InsertText: &insertText,
			Detail:     &detail,
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: overview + fmt.Sprintf("_module %s_\n\n", decl.Meta().ModuleName) + docs,
			},
		})
	}
	return &completionItems, nil
}

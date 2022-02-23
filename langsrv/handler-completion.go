package langsrv

import (
	"fmt"
	"strings"

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

	completionItems := []protocol.CompletionItem{}
	for _, decl := range globalDeclarations {
		insertText := insertTextForDecl(decl)
		var detail string
		if decl.Meta().ModuleName != "" {
			detail = string(decl.Meta().ModuleName) +
				"." +
				string(decl.DeclName())
		}
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:         string(decl.DeclName()),
			Kind:          completionItemKindForDecl(decl),
			InsertText:    &insertText,
			Detail:        &detail,
			Documentation: documentationMarkupContentForDecl(decl),
		})
	}
	return &completionItems, nil
}

func insertTextForDecl(decl ast.Decl) string {
	switch decl := decl.(type) {
	case ast.DeclFunc:
		return insertTextForCallableDeclParams(decl, decl.Impl.Parameters)
	case ast.DeclExternFunc:
		return insertTextForCallableDeclParams(decl, decl.Parameters)
	case ast.DeclData:
		return insertTextForCallableDeclFields(decl, decl.Fields)
	default:
		return string(decl.DeclName())
	}
}

func insertTextForCallableDeclParams(decl ast.Decl, parameters []ast.DeclParameter) string {
	paramNames := make([]string, len(parameters))
	for i, param := range parameters {
		paramNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, paramNames)
}

func insertTextForCallableDeclFields(decl ast.Decl, fields []ast.DeclField) string {
	fieldNames := make([]string, len(fields))
	for i, param := range fields {
		fieldNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, fieldNames)
}

func insertTextForCallableDecl(decl ast.Decl, parameters []string) string {
	if len(parameters) == 0 {
		return string(decl.DeclName())
	}
	if len(parameters) == 1 {
		return fmt.Sprintf("%s %s", decl.DeclName(), parameters[0])
	}
	return fmt.Sprintf("(%s %s)", decl.DeclName(), strings.Join(parameters, ", "))
}

func documentationMarkupContentForDecl(decl ast.Decl) *protocol.MarkupContent {
	var docs string
	if documented, ok := decl.(ast.Documented); ok {
		docs = documented.ProvidedDocs().Content + "\n\n"
	}
	var overview string
	if overviewable, ok := decl.(ast.Overviewable); ok {
		overview = "```lithia\n" + overviewable.DeclOverview() + "\n```\n\n"
	}
	return &protocol.MarkupContent{
		Kind:  protocol.MarkupKindMarkdown,
		Value: overview + fmt.Sprintf("_module %s_\n\n", decl.Meta().ModuleName) + docs,
	}
}

func completionItemKindForDecl(decl ast.Decl) *protocol.CompletionItemKind {
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
	return &kind
}

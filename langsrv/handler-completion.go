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

	completionItems := []protocol.CompletionItem{}
	for _, imported := range rc.accessibleDeclarations(context) {
		insertText := insertTextForImportedDecl(imported)
		var detail string
		if imported.decl.Meta().ModuleName != "" {
			detail = string(imported.decl.Meta().ModuleName) +
				"." +
				string(imported.decl.DeclName())
		}
		var importPrefix string
		if imported.importDecl != nil {
			importPrefix = fmt.Sprintf("%s.", imported.importDecl.DeclName())
		}
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:         importPrefix + string(imported.decl.DeclName()),
			Kind:          completionItemKindForDecl(imported.decl),
			InsertText:    &insertText,
			Detail:        &detail,
			Documentation: documentationMarkupContentForDecl(imported.decl),
		})
	}
	return &completionItems, nil
}

func insertTextForImportedDecl(imported importedDecl) string {
	var importPrefix string
	if imported.importDecl != nil {
		importPrefix = fmt.Sprintf("%s.", imported.importDecl.DeclName())
	}

	switch decl := imported.decl.(type) {
	case ast.DeclFunc:
		return insertTextForCallableDeclParams(decl, importPrefix, decl.Impl.Parameters)
	case ast.DeclExternFunc:
		return insertTextForCallableDeclParams(decl, importPrefix, decl.Parameters)
	case ast.DeclData:
		return insertTextForCallableDeclFields(decl, importPrefix, decl.Fields)
	default:
		return string(decl.DeclName())
	}
}

func insertTextForCallableDeclParams(decl ast.Decl, importPrefix string, parameters []ast.DeclParameter) string {
	paramNames := make([]string, len(parameters))
	for i, param := range parameters {
		paramNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, importPrefix, paramNames)
}

func insertTextForCallableDeclFields(decl ast.Decl, importPrefix string, fields []ast.DeclField) string {
	fieldNames := make([]string, len(fields))
	for i, param := range fields {
		fieldNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, importPrefix, fieldNames)
}

func insertTextForCallableDecl(decl ast.Decl, importPrefix string, parameters []string) string {
	if len(parameters) == 0 {
		return string(decl.DeclName())
	}
	if len(parameters) == 1 {
		return fmt.Sprintf("%s%s %s", importPrefix, decl.DeclName(), parameters[0])
	}
	return fmt.Sprintf("(%s%s %s)", importPrefix, decl.DeclName(), strings.Join(parameters, ", "))
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

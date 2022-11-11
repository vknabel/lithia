package langsrv

import (
	"fmt"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
)

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	rc := NewReqContextAtPosition(&params.TextDocumentPositionParams)

	targetNode, err := rc.findNode()
	if err != nil {
		return nil, err
	}

	completionItems := []protocol.CompletionItem{}

	switch targetNode.Type() {
	case parser.TYPE_NODE_IMPORT_DECLARATION:
		kind := protocol.CompletionItemKindModule
		mods, err := ls.resolver.ImportableModules(rc.module.Package())
		if err != nil {
			return nil, err
		}
		for _, mod := range mods {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: string(mod.RelativeName),
				Kind:  &kind,
				Documentation: &protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: fmt.Sprintf("`import %s`", mod.RelativeName),
				},
			})
		}
		return completionItems, nil
	case parser.TYPE_NODE_MEMBER_ACCESS, ".":
		return rc.textDocumentMemberAccessCompletionItems(context)
	}

	for _, imported := range rc.accessibleDeclarations(context) {
		completion := rc.generalCompletionItemsForDecl(imported)
		completionItems = append(completionItems, completion)
	}
	return &completionItems, nil
}

func (rc *ReqContext) textDocumentMemberAccessCompletionItems(context *glsp.Context) ([]protocol.CompletionItem, error) {
	defaultScope := rc.accessibleDeclarations(context)
	targetNode, err := rc.findNode()
	if err != nil {
		return nil, err
	}
	accessedNode := targetNode
	if accessedNode.Type() == "." {
		accessTarget := targetNode.PrevNamedSibling()
		if accessTarget != nil {
			accessedNode = accessTarget
		}
	}
	accessedExpr := accessedNode.Content([]byte(rc.textDocumentEntry.item.Text))
	for _, imported := range defaultScope {
		switch decl := imported.decl.(type) {
		case ast.DeclModule:
			if string(decl.DeclName()) != accessedExpr {
				continue
			}
			moduleDecls := rc.moduleDeclarations()
			scope := make([]importedDecl, 0)
			for _, moduleDecl := range moduleDecls {
				if !moduleDecl.IsExportedDecl() {
					continue
				}
				scope = append(scope, importedDecl{decl: moduleDecl, module: rc.textDocumentEntry.module, importDecl: nil})
			}

			var completionItems []protocol.CompletionItem
			for _, imported := range scope {
				completionItems = append(
					completionItems,
					rc.generalCompletionItemsForDecl(imported),
				)
			}
			return completionItems, nil
		}
	}
	var completionItems []protocol.CompletionItem
	for _, imported := range defaultScope {
		switch imported.decl.(type) {
		case ast.DeclData, ast.DeclExternType:
			completionItems = append(completionItems, rc.memberAccessCompletionItemsForDecl(imported, accessedExpr)...)
		}
	}
	return completionItems, nil
}

func (rc *ReqContext) generalCompletionItemsForDecl(imported importedDecl) protocol.CompletionItem {
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
	return protocol.CompletionItem{
		Label:         importPrefix + string(imported.decl.DeclName()),
		Kind:          completionItemKindForDecl(imported.decl),
		InsertText:    &insertText,
		Detail:        &detail,
		Documentation: documentationMarkupContentForDecl(imported.decl),
	}
}

func (rc *ReqContext) memberAccessCompletionItemsForDecl(imported importedDecl, accessedExpr string) []protocol.CompletionItem {
	switch decl := imported.decl.(type) {
	case ast.DeclModule:
		moduleDecls := rc.moduleDeclarations()
		importedDecls := make([]importedDecl, 0, len(moduleDecls))
		for _, moduleDecl := range moduleDecls {
			if !moduleDecl.IsExportedDecl() {
				continue
			}
			importedDecls = append(importedDecls, importedDecl{decl: moduleDecl, module: rc.textDocumentEntry.module, importDecl: nil})
		}
		completionItems := make([]protocol.CompletionItem, 0)
		for _, current := range importedDecls {
			completion := rc.generalCompletionItemsForDecl(current)
			completionItems = append(completionItems, completion)
		}
		return completionItems
	case ast.DeclImport:
		// TODO
		return nil
	case ast.DeclImportMember:
		// TODO
		return nil
	case ast.DeclExternType:
		completionItems := make([]protocol.CompletionItem, 0, len(decl.Fields))
		for _, field := range decl.Fields {
			insertText := string(field.DeclName())
			var detail string
			if imported.decl.Meta().ModuleName != "" {
				detail = string(imported.decl.Meta().ModuleName) +
					"." +
					string(imported.decl.DeclName()) +
					"." +
					string(field.DeclName())
			}

			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            string(field.DeclName()),
				Kind:             completionItemKindForDecl(imported.decl),
				InsertText:       &insertText,
				Detail:           &detail,
				CommitCharacters: []string{"."},
				Documentation:    documentationMarkupContentForDecl(imported.decl),
			})
		}
		return completionItems
	case ast.DeclData:
		completionItems := make([]protocol.CompletionItem, 0, len(decl.Fields))
		for _, field := range decl.Fields {
			insertText := string(field.DeclName())
			var detail string
			if imported.decl.Meta().ModuleName != "" {
				detail = string(imported.decl.Meta().ModuleName) +
					"." +
					string(imported.decl.DeclName()) +
					"." +
					string(field.DeclName())
			}

			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            string(field.DeclName()),
				Kind:             completionItemKindForDecl(imported.decl),
				InsertText:       &insertText,
				Detail:           &detail,
				CommitCharacters: []string{"."},
				Documentation:    documentationMarkupContentForDecl(imported.decl),
			})
		}
		return completionItems
	default:
		return nil
	}
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

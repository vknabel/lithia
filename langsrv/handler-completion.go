package langsrv

import (
	"fmt"
	"unicode"

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

	insertAsStatement := false
	contextReferenceType := targetNode.Type()
	if contextReferenceType == parser.TYPE_NODE_IDENTIFIER {
		contextReferenceType = targetNode.Parent().Type()
	}
	switch contextReferenceType {
	case parser.TYPE_NODE_FUNCTION_DECLARATION, "function_body", parser.TYPE_NODE_FUNCTION_LITERAL,
		parser.TYPE_NODE_SOURCE_FILE:
		insertAsStatement = true
	case parser.TYPE_NODE_ARRAY_LITERAL,
		parser.TYPE_NODE_BINARY_EXPRESSION,
		parser.TYPE_NODE_DICT_LITERAL, parser.TYPE_NODE_DICT_ENTRY:
		insertAsStatement = false
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
		return rc.textDocumentMemberAccessCompletionItems(context, insertAsStatement)
	}

	for _, imported := range rc.accessibleDeclarations(context) {
		completions := rc.generalCompletionItemsForDecl(imported, insertAsStatement)
		completionItems = append(completionItems, completions...)
	}
	completionItems = append(completionItems, rc.keywordCompletionItems(insertAsStatement)...)
	return completionItems, nil
}

func (rc *ReqContext) keywordCompletionItems(asStatement bool) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	keywordKind := protocol.CompletionItemKindKeyword
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "type ${1:Type} {\n    $0\n}"
		detail := "type expression"
		item := protocol.CompletionItem{
			Label:            "type",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	if !asStatement {
		return completionItems
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "let ${1:var} = $0"
		detail := "constant declaration"
		item := protocol.CompletionItem{
			Label:            "let",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		name := string(rc.module.RelativeName)
		insertText := fmt.Sprintf("module ${0:%s}", name)
		detail := "module declaration"
		item := protocol.CompletionItem{
			Label:            "module",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "import $0"
		detail := "import declaration"
		item := protocol.CompletionItem{
			Label:            "import",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "import ${1:alias} = ${2:module.path}"
		detail := "import declaration"
		item := protocol.CompletionItem{
			Label:            "import =",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "func ${1:name} { $2=>\n    $0\n}"
		detail := "function declaration"
		item := protocol.CompletionItem{
			Label:            "func",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "enum ${1:Name} {\n    $0\n}"
		detail := "enum declaration"
		item := protocol.CompletionItem{
			Label:            "enum",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "data ${1:Name}"
		detail := "data declaration"
		item := protocol.CompletionItem{
			Label:            "data",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	{
		insertFormat := protocol.InsertTextFormatSnippet
		insertText := "extern ${0:Name}"
		detail := "extern declaration"
		item := protocol.CompletionItem{
			Label:            "extern",
			Kind:             &keywordKind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
		}
		completionItems = append(completionItems, item)
	}
	return completionItems
}

func (rc *ReqContext) textDocumentMemberAccessCompletionItems(context *glsp.Context, asStatement bool) ([]protocol.CompletionItem, error) {
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
					rc.generalCompletionItemsForDecl(imported, asStatement)...,
				)
			}
			return completionItems, nil
		}
	}
	var completionItems []protocol.CompletionItem
	for _, imported := range defaultScope {
		switch imported.decl.(type) {
		case ast.DeclData, ast.DeclExternType:
			completionItems = append(completionItems, rc.memberAccessCompletionItemsForDecl(imported, accessedExpr, asStatement)...)
		}
	}
	return completionItems, nil
}

func (rc *ReqContext) generalCompletionItemsForDecl(imported importedDecl, asStatement bool) []protocol.CompletionItem {
	insertFormat, insertText := insertTextForImportedDecl(imported, asStatement)
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
	completionItems := make([]protocol.CompletionItem, 0, 2)
	declCompletion := protocol.CompletionItem{
		Label:            importPrefix + string(imported.decl.DeclName()),
		Kind:             completionItemKindForDecl(imported.decl),
		InsertText:       &insertText,
		InsertTextFormat: &insertFormat,
		Detail:           &detail,
		Documentation:    documentationMarkupContentForDecl(imported.decl),
	}
	completionItems = append(completionItems, declCompletion)
	if enumDecl, ok := imported.decl.(ast.DeclEnum); ok {
		kind := protocol.CompletionItemKindSnippet
		insertText := fmt.Sprintf("type %s%s {", importPrefix, string(enumDecl.DeclName()))
		if len(enumDecl.Cases) != 0 {
			insertText += "\n"
		}
		for i, enumCase := range enumDecl.Cases {
			caseName := string(enumCase.DeclName())
			lowerHead := unicode.ToLower(rune(caseName[0]))
			varName := string(lowerHead) + caseName[1:]
			insertText += fmt.Sprintf("    %s: { %s => $%d },\n", caseName, varName, i+1)
		}
		insertText += "}"
		insertFormat := protocol.InsertTextFormatSnippet
		typeCompletion := protocol.CompletionItem{
			Label:            "type " + importPrefix + string(enumDecl.DeclName()),
			Kind:             &kind,
			InsertText:       &insertText,
			InsertTextFormat: &insertFormat,
			Detail:           &detail,
			Documentation:    documentationMarkupContentForDecl(enumDecl),
		}
		completionItems = append(completionItems, typeCompletion)
	}
	return completionItems
}

func (rc *ReqContext) memberAccessCompletionItemsForDecl(imported importedDecl, accessedExpr string, asStatement bool) []protocol.CompletionItem {
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
			completion := rc.generalCompletionItemsForDecl(current, asStatement)
			completionItems = append(completionItems, completion...)
		}
		return completionItems
	case ast.DeclImport:
		return nil
	case ast.DeclImportMember:
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

func insertTextForImportedDecl(imported importedDecl, asStatement bool) (protocol.InsertTextFormat, string) {
	var importPrefix string
	if imported.importDecl != nil {
		importPrefix = fmt.Sprintf("%s.", imported.importDecl.DeclName())
	}

	switch decl := imported.decl.(type) {
	case ast.DeclFunc:
		return insertTextForCallableDeclParams(decl, importPrefix, decl.Impl.Parameters, asStatement)
	case ast.DeclExternFunc:
		return insertTextForCallableDeclParams(decl, importPrefix, decl.Parameters, asStatement)
	case ast.DeclData:
		return insertTextForCallableDeclFields(decl, importPrefix, decl.Fields, asStatement)
	default:
		return protocol.InsertTextFormatPlainText, string(decl.DeclName())
	}
}

func insertTextForCallableDeclParams(decl ast.Decl, importPrefix string, parameters []ast.DeclParameter, asStatement bool) (protocol.InsertTextFormat, string) {
	paramNames := make([]string, len(parameters))
	for i, param := range parameters {
		paramNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, importPrefix, paramNames, asStatement)
}

func insertTextForCallableDeclFields(decl ast.Decl, importPrefix string, fields []ast.DeclField, asStatement bool) (protocol.InsertTextFormat, string) {
	fieldNames := make([]string, len(fields))
	for i, param := range fields {
		fieldNames[i] = string(param.DeclName())
	}
	return insertTextForCallableDecl(decl, importPrefix, fieldNames, asStatement)
}

func insertTextForCallableDecl(decl ast.Decl, importPrefix string, parameters []string, asStatement bool) (protocol.InsertTextFormat, string) {
	if len(parameters) == 0 {
		return protocol.InsertTextFormatPlainText, string(decl.DeclName())
	}
	insertText := fmt.Sprintf("%s%s", importPrefix, decl.DeclName())
	for i, parameter := range parameters {
		if i > 0 {
			insertText += ", "
		} else {
			insertText += " "
		}
		insertText += fmt.Sprintf("${%d:%s}", i+1, parameter)
	}
	if len(parameters) == 1 || asStatement {
		return protocol.InsertTextFormatSnippet, insertText
	}
	return protocol.InsertTextFormatSnippet, fmt.Sprintf("(%s)", insertText)
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

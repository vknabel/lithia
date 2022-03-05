package langsrv

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/resolution"
)

type ReqContext struct {
	textDocument protocol.TextDocumentIdentifier
	position     protocol.Position
	textDocumentEntry
}

func NewReqContext(textDocument protocol.TextDocumentIdentifier) *ReqContext {
	return &ReqContext{
		textDocument:      textDocument,
		textDocumentEntry: *ls.documentCache.documents[textDocument.URI],
	}
}

func NewReqContextAtPosition(position *protocol.TextDocumentPositionParams) *ReqContext {
	return &ReqContext{
		textDocument:      position.TextDocument,
		position:          position.Position,
		textDocumentEntry: *ls.documentCache.documents[position.TextDocument.URI],
	}
}

func (rc *ReqContext) findToken() (string, *protocol.Range, error) {
	node, err := rc.findNode()
	if err != nil || node == nil {
		return "", nil, err
	}
	contents := rc.textDocumentEntry.item.Text
	name := node.Content([]byte(contents))
	return name, &protocol.Range{
		Start: protocol.Position{
			Line:      uint32(node.StartPoint().Row),
			Character: uint32(node.StartPoint().Column),
		},
		End: protocol.Position{
			Line:      uint32(node.EndPoint().Row),
			Character: uint32(node.EndPoint().Column),
		},
	}, nil
}

func (rc *ReqContext) findNode() (*sitter.Node, error) {
	node := NodeAtPosition(rc.fileParser.Tree.RootNode(), rc.position)
	return node, nil
}

func (rc *ReqContext) globalDeclarations(context *glsp.Context) []importedDecl {
	if rc.sourceFile == nil {
		return nil
	}

	globals := make([]importedDecl, 0)
	for _, moduleDecl := range rc.moduleDeclarations() {
		globals = append(globals, importedDecl{decl: moduleDecl, module: rc.textDocumentEntry.module, importDecl: nil})
	}
	globals = append(globals, rc.importedDeclarations(context)...)

	return globals
}

func (rc *ReqContext) moduleDeclarations() []ast.Decl {
	if rc.sourceFile == nil {
		return nil
	}
	globalDeclarations := rc.sourceFile.Declarations
	for _, sameModuleFile := range rc.textDocumentEntry.module.Files {
		fileUrl := "file://" + sameModuleFile
		if rc.item.URI == fileUrl {
			continue
		}
		docEntry := ls.documentCache.documents[fileUrl]
		if docEntry == nil || docEntry.sourceFile == nil {
			continue
		}

		globalDeclarations = append(globalDeclarations, docEntry.sourceFile.ExportedDeclarations()...)
	}
	return globalDeclarations
}

type importedDecl struct {
	decl       ast.Decl
	module     resolution.ResolvedModule
	importDecl *ast.DeclImport
}

func (rc *ReqContext) importedDeclarations(context *glsp.Context) []importedDecl {
	if rc.sourceFile == nil {
		return nil
	}

	globals := make([]importedDecl, 0)

	resolvedPrelude, err := ls.resolver.ResolveModuleFromPackage(rc.textDocumentEntry.module.Package(), "prelude")
	if err != nil {
		ls.server.Log.Error(err.Error())
	} else {
		openModuleTextDocumentsIfNeeded(context, resolvedPrelude)
	}

	for _, sameModuleFile := range resolvedPrelude.Files {
		fileUri := "file://" + sameModuleFile
		if ls.documentCache.documents[fileUri] == nil {
			continue
		}
		entry := ls.documentCache.documents[fileUri]
		if entry.sourceFile == nil {
			continue
		}
		for _, decl := range entry.sourceFile.ExportedDeclarations() {
			globals = append(globals, importedDecl{decl, resolvedPrelude, nil})
		}
	}

	for _, decl := range rc.sourceFile.Declarations {
		if _, ok := decl.(ast.DeclImport); !ok {
			continue
		}
		importDecl := decl.(ast.DeclImport)
		resolvedModule, err := ls.resolver.ResolveModuleFromPackage(rc.textDocumentEntry.module.Package(), importDecl.ModuleName)
		if err != nil {
			ls.server.Log.Error(err.Error())
		} else {
			openModuleTextDocumentsIfNeeded(context, resolvedModule)
		}

		for _, sameModuleFile := range resolvedModule.Files {
			fileUri := "file://" + sameModuleFile
			if ls.documentCache.documents[fileUri] == nil {
				continue
			}
			entry := ls.documentCache.documents[fileUri]
			if entry.sourceFile == nil {
				continue
			}
			for _, decl := range entry.sourceFile.ExportedDeclarations() {
				globals = append(globals, importedDecl{decl, resolvedModule, &importDecl})
			}
		}
	}
	return globals
}

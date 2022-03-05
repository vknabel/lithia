package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	mod := langserver.resolver.ResolvePackageAndModuleForReferenceFile(params.TextDocument.URI)
	entry := langserver.documentCache.documents[params.TextDocument.URI]
	text := entry.item.Text
	for _, event := range params.ContentChanges {
		switch e := event.(type) {
		case protocol.TextDocumentContentChangeEvent:
			text = text[:e.Range.Start.IndexIn(text)] + e.Text + text[e.Range.End.IndexIn(text):]
		case protocol.TextDocumentContentChangeEventWhole:
			text = e.Text
		}
	}
	entry.item.Text = text
	syntaxErrs := make([]parser.SyntaxError, 0)
	fileParser, errs := entry.parser.Parse(mod.AbsoluteModuleName(), string(params.TextDocument.URI), text)
	syntaxErrs = append(syntaxErrs, errs...)
	sourceFile, errs := fileParser.ParseSourceFile()
	syntaxErrs = append(syntaxErrs, errs...)
	langserver.documentCache.documents[params.TextDocument.URI].fileParser = fileParser
	langserver.documentCache.documents[params.TextDocument.URI].sourceFile = sourceFile

	analyzeErrs := analyzeErrorsForSourceFile(context, mod, *sourceFile)
	publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), syntaxErrs, analyzeErrs)
	return nil
}

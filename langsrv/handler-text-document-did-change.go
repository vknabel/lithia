package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	entry := langserver.documentCache.documents[params.TextDocument.URI]
	text := entry.item.Text
	for _, event := range params.ContentChanges {
		switch e := event.(type) {
		case protocol.TextDocumentContentChangeEvent:
			langserver.server.Log.Infof("from: %s", text)
			text = text[:e.Range.Start.IndexIn(text)] + e.Text + text[e.Range.End.IndexIn(text):]
			langserver.server.Log.Infof("to: %s", text)
		case protocol.TextDocumentContentChangeEventWhole:
			text = e.Text
		}
	}
	entry.item.Text = text
	fileParser, errs := entry.parser.Parse("default-module", string(params.TextDocument.URI), text)
	if len(errs) > 0 {
		// TODO: syntax errors
		return parser.NewGroupedSyntaxError(errs)
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if len(errs) > 0 {
		// TODO: syntax errors
		return parser.NewGroupedSyntaxError(errs)
	}
	langserver.server.Log.Infof("%s: %s", params.TextDocument.URI, text)
	langserver.documentCache.documents[params.TextDocument.URI].fileParser = fileParser
	langserver.documentCache.documents[params.TextDocument.URI].sourceFile = sourceFile
	return nil
}

package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	lithiaParser := parser.NewParser()
	fileParser, errs := lithiaParser.Parse("default-module", string(params.TextDocument.URI), params.TextDocument.Text)
	if len(errs) > 0 {
		publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), errs)
		return nil
	}
	sourceFile, errs := fileParser.ParseSourceFile()
	if len(errs) > 0 {
		publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), errs)
		return nil
	}
	langserver.documentCache.documents[params.TextDocument.URI] = &textDocumentEntry{
		item:       params.TextDocument,
		parser:     lithiaParser,
		fileParser: fileParser,
		sourceFile: sourceFile,
	}
	publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), nil)
	return nil
}

package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	lithiaParser := parser.NewParser()

	syntaxErrs := make([]parser.SyntaxError, 0)
	fileParser, errs := lithiaParser.Parse("default-module", string(params.TextDocument.URI), params.TextDocument.Text)
	syntaxErrs = append(syntaxErrs, errs...)
	sourceFile, errs := fileParser.ParseSourceFile()
	syntaxErrs = append(syntaxErrs, errs...)
	langserver.documentCache.documents[params.TextDocument.URI] = &textDocumentEntry{
		item:       params.TextDocument,
		parser:     lithiaParser,
		fileParser: fileParser,
		sourceFile: sourceFile,
	}
	publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), syntaxErrs)
	return nil
}

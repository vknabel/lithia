package langsrv

import (
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/resolution"
)

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	lithiaParser := parser.NewParser()
	mod := langserver.resolver.ResolvePackageAndModuleForReferenceFile(params.TextDocument.URI)
	openModuleTextDocumentsIfNeeded(context, mod)

	syntaxErrs := make([]parser.SyntaxError, 0)
	fileParser, errs := lithiaParser.Parse(mod.AbsoluteModuleName(), string(params.TextDocument.URI), params.TextDocument.Text)
	syntaxErrs = append(syntaxErrs, errs...)
	sourceFile, errs := fileParser.ParseSourceFile()
	syntaxErrs = append(syntaxErrs, errs...)
	langserver.documentCache.documents[params.TextDocument.URI] = &textDocumentEntry{
		item:       params.TextDocument,
		parser:     lithiaParser,
		fileParser: fileParser,
		sourceFile: sourceFile,
		module:     mod,
	}

	publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), syntaxErrs)
	return nil
}

func openModuleTextDocumentsIfNeeded(context *glsp.Context, mod resolution.ResolvedModule) {
	fileNames := mod.Files
	if mod.Package().Manifest != nil {
		fileNames = append(fileNames, mod.Package().Manifest.Path)
	}
	for _, filePath := range fileNames {
		fileUri := "file://" + filePath
		if langserver.documentCache.documents[fileUri] != nil {
			continue
		}
		lithiaParser := parser.NewParser()

		syntaxErrs := make([]parser.SyntaxError, 0)
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		fileParser, errs := lithiaParser.Parse(mod.AbsoluteModuleName(), string(fileUri), string(bytes))
		syntaxErrs = append(syntaxErrs, errs...)
		sourceFile, errs := fileParser.ParseSourceFile()
		syntaxErrs = append(syntaxErrs, errs...)

		langserver.documentCache.documents[fileUri] = &textDocumentEntry{
			parser:     lithiaParser,
			fileParser: fileParser,
			sourceFile: sourceFile,
			module:     mod,
		}

		publishSyntaxErrorDiagnosticsForFile(context, fileUri, syntaxErrs)
	}
}

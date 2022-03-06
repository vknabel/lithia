package langsrv

import (
	"fmt"
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
	"github.com/vknabel/lithia/resolution"
)

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	lithiaParser := parser.NewParser()
	mod := ls.resolver.ResolvePackageAndModuleForReferenceFile(params.TextDocument.URI)
	openModuleTextDocumentsIfNeeded(context, mod)

	syntaxErrs := make([]parser.SyntaxError, 0)
	fileParser, errs := lithiaParser.Parse(mod.AbsoluteModuleName(), string(params.TextDocument.URI), params.TextDocument.Text)
	syntaxErrs = append(syntaxErrs, errs...)
	sourceFile, errs := fileParser.ParseSourceFile()
	syntaxErrs = append(syntaxErrs, errs...)
	ls.documentCache.documents[params.TextDocument.URI] = &textDocumentEntry{
		item:       params.TextDocument,
		parser:     lithiaParser,
		fileParser: fileParser,
		sourceFile: sourceFile,
		module:     mod,
	}

	analyzeErrs := analyzeErrorsForSourceFile(context, mod, *sourceFile)
	publishSyntaxErrorDiagnostics(context, params.TextDocument.URI, uint32(params.TextDocument.Version), syntaxErrs, analyzeErrs)
	return nil
}

func openModuleTextDocumentsIfNeeded(context *glsp.Context, mod resolution.ResolvedModule) {
	fileNames := mod.Files
	if mod.Package().Manifest != nil {
		fileNames = append(fileNames, mod.Package().Manifest.Path)
	}
	for _, filePath := range fileNames {
		fileUri := "file://" + filePath
		if ls.documentCache.documents[fileUri] != nil {
			continue
		}
		lithiaParser := parser.NewParser()

		syntaxErrs := make([]parser.SyntaxError, 0)
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			ls.server.Log.Errorf("failed to read %s, due %s", fileUri, err.Error())
			continue
		}
		contents := string(bytes)
		fileParser, errs := lithiaParser.Parse(mod.AbsoluteModuleName(), string(fileUri), contents)
		syntaxErrs = append(syntaxErrs, errs...)
		sourceFile, errs := fileParser.ParseSourceFile()
		syntaxErrs = append(syntaxErrs, errs...)

		ls.documentCache.documents[fileUri] = &textDocumentEntry{
			parser:     lithiaParser,
			fileParser: fileParser,
			sourceFile: sourceFile,
			module:     mod,
		}

		analyzeErrs := analyzeErrorsForSourceFile(context, mod, *sourceFile)
		publishSyntaxErrorDiagnosticsForFile(context, fileUri, syntaxErrs, analyzeErrs)
	}
}

func analyzeErrorsForSourceFile(context *glsp.Context, mod resolution.ResolvedModule, sourceFile ast.SourceFile) []analyzeError {
	analyzeErrs := make([]analyzeError, 0)
	for _, decl := range sourceFile.Declarations {
		if _, ok := decl.(ast.DeclImport); !ok {
			continue
		}
		importDecl := decl.(ast.DeclImport)
		resolvedModule, err := ls.resolver.ResolveModuleFromPackage(mod.Package(), importDecl.ModuleName)
		if err != nil {
			analyzeErrs = append(
				analyzeErrs,
				newAnalyzeErrorAtLocation("error", fmt.Sprintf("module %s not found", importDecl.ModuleName), importDecl.Meta().Source),
			)
		} else {
			openModuleTextDocumentsIfNeeded(context, resolvedModule)
		}
	}
	return analyzeErrs
}

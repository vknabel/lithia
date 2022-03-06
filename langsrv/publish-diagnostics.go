package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
)

type analyzeError struct {
	kind    string
	message string
	source  *ast.Source
}

func newAnalyzeErrorAtLocation(kind string, message string, source *ast.Source) analyzeError {
	return analyzeError{
		kind:    kind,
		message: message,
		source:  source,
	}
}

func publishSyntaxErrorDiagnostics(context *glsp.Context, textDocumentURI protocol.URI, version uint32, errs []parser.SyntaxError, analyzeErrs []analyzeError) {
	diagnostics := make([]protocol.Diagnostic, len(errs)+len(analyzeErrs))
	for i, err := range errs {
		diagnostics[i] = syntaxErrorToDiagnostic(err)
	}
	for i, err := range analyzeErrs {
		diagnostics[len(errs)+i] = analyzeErrorToDiagnostic(err)
	}

	versionRef := version
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         textDocumentURI,
		Version:     &versionRef,
		Diagnostics: diagnostics,
	})
}

func publishSyntaxErrorDiagnosticsForFile(context *glsp.Context, textDocumentURI protocol.URI, errs []parser.SyntaxError, analyzeErrs []analyzeError) {
	diagnostics := make([]protocol.Diagnostic, len(errs)+len(analyzeErrs))
	for i, err := range errs {
		diagnostics[i] = syntaxErrorToDiagnostic(err)
	}
	for i, err := range analyzeErrs {
		diagnostics[len(errs)+i] = analyzeErrorToDiagnostic(err)
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         textDocumentURI,
		Version:     nil,
		Diagnostics: diagnostics,
	})
}

func syntaxErrorToDiagnostic(err parser.SyntaxError) protocol.Diagnostic {
	return protocol.Diagnostic{
		Source:   &lsName,
		Range:    rangeFromParserSourceLocation(err.SourceLocation),
		Severity: newSeverityRef(protocol.DiagnosticSeverityError),
		Message:  err.Message,
	}
}

func analyzeErrorToDiagnostic(err analyzeError) protocol.Diagnostic {
	return protocol.Diagnostic{
		Source:   &lsName,
		Range:    rangeFromAstSourceLocation(err.source),
		Severity: newSeverityRef(protocol.DiagnosticSeverityError),
		Message:  err.message,
	}
}

func rangeFromParserSourceLocation(location parser.SourceLocation) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      location.Node.StartPoint().Row,
			Character: location.Node.StartPoint().Column,
		},
		End: protocol.Position{
			Line:      location.Node.EndPoint().Row,
			Character: location.Node.EndPoint().Column,
		},
	}
}

func rangeFromAstSourceLocation(location *ast.Source) protocol.Range {
	if location == nil {
		return protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 0},
		}
	}
	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(location.Start.Line),
			Character: uint32(location.Start.Column),
		},
		End: protocol.Position{
			Line:      uint32(location.End.Line),
			Character: uint32(location.End.Column),
		},
	}
}

func newSeverityRef(sev protocol.DiagnosticSeverity) *protocol.DiagnosticSeverity {
	severity := sev
	return &severity
}

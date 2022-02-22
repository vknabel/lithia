package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/parser"
)

func publishSyntaxErrorDiagnostics(context *glsp.Context, textDocumentURI protocol.URI, version uint32, errs []parser.SyntaxError) {
	diagnostics := make([]protocol.Diagnostic, len(errs))
	for i, err := range errs {
		diagnostics[i] = syntaxErrorToDiagnostic(err)
	}

	versionRef := version
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         textDocumentURI,
		Version:     &versionRef,
		Diagnostics: diagnostics,
	})
}

func publishSyntaxErrorDiagnosticsForFile(context *glsp.Context, textDocumentURI protocol.URI, errs []parser.SyntaxError) {
	diagnostics := make([]protocol.Diagnostic, len(errs))
	for i, err := range errs {
		diagnostics[i] = syntaxErrorToDiagnostic(err)
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
		Range:    rangeFromSourceLocation(err.SourceLocation),
		Severity: newSeverityRef(protocol.DiagnosticSeverityError),
		Message:  err.Message,
	}
}

func rangeFromSourceLocation(location parser.SourceLocation) protocol.Range {
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

func newSeverityRef(sev protocol.DiagnosticSeverity) *protocol.DiagnosticSeverity {
	severity := sev
	return &severity
}

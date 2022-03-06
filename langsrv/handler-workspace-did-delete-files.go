package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func workspaceDidDeleteFiles(context *glsp.Context, params *protocol.DeleteFilesParams) error {
	for _, deleted := range params.Files {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         deleted.URI,
			Version:     nil,
			Diagnostics: []protocol.Diagnostic{},
		})
		delete(ls.documentCache.documents, deleted.URI)
	}
	return nil
}

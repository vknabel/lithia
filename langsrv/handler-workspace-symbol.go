package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func workspaceSymbol(context *glsp.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	symbols := make([]protocol.SymbolInformation, 0)

	for uri, document := range ls.documentCache.documents {
		exported := document.sourceFile.ExportedDeclarations()
		for _, decl := range exported {
			containerName := string(decl.Meta().ModuleName)
			symbols = append(symbols, protocol.SymbolInformation{
				Name:          string(decl.DeclName()),
				ContainerName: &containerName,
				Kind:          symbolKindClassForDecl(decl),
				Location: protocol.Location{
					URI:   uri,
					Range: rangeFromAstSourceLocation(decl.Meta().Source),
				},
			})
		}
	}
	return symbols, nil
}

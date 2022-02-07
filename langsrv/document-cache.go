package langsrv

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
)

type documentCache struct {
	documents map[protocol.URI]*textDocumentEntry
}

type textDocumentEntry struct {
	item       protocol.TextDocumentItem
	parser     *parser.Parser
	fileParser *parser.FileParser
	sourceFile *ast.SourceFile
}

package langsrv

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"

	_ "github.com/tliron/kutil/logging/simple"
)

var lsName = "lithia"
var debug = true
var handler protocol.Handler

type lithiaLangserver struct {
	server        *server.Server
	documentCache *documentCache
}

var langserver lithiaLangserver = lithiaLangserver{
	documentCache: &documentCache{documents: make(map[protocol.URI]*textDocumentEntry)},
}

func init() {
	logging.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,

		TextDocumentDidOpen:   textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,

		TextDocumentHover:          textDocumentHover,
		TextDocumentCompletion:     textDocumentCompletion,
		TextDocumentDefinition:     textDocumentDefinition,
		TextDocumentTypeDefinition: textDocumentTypeDefinition,
		TextDocumentDeclaration:    textDocumentDeclaration,

		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
	}
}

func RunStdio() error {
	langserver.server = server.NewServer(&handler, lsName, debug)
	return langserver.server.RunStdio()
}

func RunIPC() error {
	langserver.server = server.NewServer(&handler, lsName, debug)
	return langserver.server.RunNodeJs()
}

func RunSocket(address string) error {
	langserver.server = server.NewServer(&handler, lsName, debug)
	return langserver.server.RunWebSocket(address)
}

func RunTCP(address string) error {
	langserver.server = server.NewServer(&handler, lsName, debug)
	return langserver.server.RunTCP(address)
}

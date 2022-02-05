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

func init() {
	logging.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,

		TextDocumentHover:          textDocumentHover,
		TextDocumentCompletion:     textDocumentCompletion,
		TextDocumentDefinition:     textDocumentDefinition,
		TextDocumentTypeDefinition: textDocumentTypeDefinition,
		TextDocumentDeclaration:    textDocumentDeclaration,

		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
	}
}

func RunStdio() error {
	server := server.NewServer(&handler, lsName, debug)
	return server.RunStdio()
}

func RunIPC() error {
	server := server.NewServer(&handler, lsName, debug)
	return server.RunNodeJs()
}

func RunSocket(address string) error {
	server := server.NewServer(&handler, lsName, debug)
	return server.RunWebSocket(address)
}

func RunTCP(address string) error {
	server := server.NewServer(&handler, lsName, debug)
	return server.RunTCP(address)
}

package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"
	"github.com/vknabel/lithia/info"
	"github.com/vknabel/lithia/resolution"

	_ "github.com/tliron/kutil/logging/simple"
)

var lsName = "lithia"
var debug = info.Debug
var handler protocol.Handler

type lithiaLangserver struct {
	resolver resolution.ModuleResolver

	server         *server.Server
	documentCache  *documentCache
	workspaceRoots []string
}

var ls lithiaLangserver = lithiaLangserver{
	resolver:       resolution.NewDefaultModuleResolver(),
	documentCache:  &documentCache{documents: make(map[protocol.URI]*textDocumentEntry)},
	workspaceRoots: []string{},
}

func init() {
	logging.Configure(2, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,

		CancelRequest: func(context *glsp.Context, params *protocol.CancelParams) error { return nil },

		TextDocumentDidOpen:   textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,
		TextDocumentDidClose:  func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error { return nil },

		WorkspaceDidDeleteFiles: workspaceDidDeleteFiles,
		WorkspaceSymbol:         workspaceSymbol,

		TextDocumentHover:          textDocumentHover,
		TextDocumentCompletion:     textDocumentCompletion,
		TextDocumentDefinition:     textDocumentDefinition,
		TextDocumentTypeDefinition: textDocumentTypeDefinition,
		TextDocumentDeclaration:    textDocumentDeclaration,
		TextDocumentDocumentSymbol: textDocumentDocumentSymbol,

		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
	}
}

func RunStdio() error {
	ls.server = server.NewServer(&handler, lsName, debug)
	return ls.server.RunStdio()
}

func RunIPC() error {
	ls.server = server.NewServer(&handler, lsName, debug)
	return ls.server.RunNodeJs()
}

func RunSocket(address string) error {
	ls.server = server.NewServer(&handler, lsName, debug)
	return ls.server.RunWebSocket(address)
}

func RunTCP(address string) error {
	ls.server = server.NewServer(&handler, lsName, debug)
	return ls.server.RunTCP(address)
}

func (ls *lithiaLangserver) setWorkspaceRoots(roots ...string) {
	ls.workspaceRoots = roots
}

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

var langserver lithiaLangserver = lithiaLangserver{
	resolver:       resolution.DefaultModuleResolver(),
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

		TextDocumentDidOpen:   textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,
		TextDocumentDidClose:  func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error { return nil },

		WorkspaceDidDeleteFiles: workspaceDidDeleteFiles,

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

func (ls *lithiaLangserver) setWorkspaceRoots(roots ...string) {
	ls.workspaceRoots = roots
}

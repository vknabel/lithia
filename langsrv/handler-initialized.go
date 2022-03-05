package langsrv

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	for _, root := range ls.workspaceRoots {
		matches, err := filepath.Glob(path.Join(strings.TrimPrefix("file://", root), "*/*.lithia"))
		if err != nil {
			ls.server.Log.Errorf("package detection failed, due %s", err)
			continue
		}
		for _, match := range matches {
			mod := ls.resolver.ResolvePackageAndModuleForReferenceFile(match)
			openModuleTextDocumentsIfNeeded(context, mod)
		}
	}
	return nil
}

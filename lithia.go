package lithia

import (
	extdocs "github.com/vknabel/lithia/external/docs"
	extfs "github.com/vknabel/lithia/external/fs"
	extos "github.com/vknabel/lithia/external/os"
	extrx "github.com/vknabel/lithia/external/rx"
	exttea "github.com/vknabel/lithia/external/tea"
	"github.com/vknabel/lithia/runtime"
)

func NewDefaultInterpreter(referenceFile string, importRoots ...string) *runtime.Interpreter {
	inter := runtime.NewIsolatedInterpreter(referenceFile, importRoots...)
	inter.ExternalDefinitions["os"] = extos.ExternalOS{}
	inter.ExternalDefinitions["rx"] = extrx.ExternalRx{}
	inter.ExternalDefinitions["docs"] = extdocs.ExternalDocs{}
	inter.ExternalDefinitions["fs"] = extfs.ExternalFS{}
	inter.ExternalDefinitions["tea"] = exttea.ExternalTea{}
	return inter
}

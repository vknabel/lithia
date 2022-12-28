package lithia

import (
	extdocs "github.com/vknabel/lithia/external/docs"
	extfs "github.com/vknabel/lithia/external/fs"
	extos "github.com/vknabel/lithia/external/os"
	extrx "github.com/vknabel/lithia/external/rx"
	"github.com/vknabel/lithia/runtime"
)

func NewDefaultInterpreter(referenceFile string, importRoots ...string) *runtime.Interpreter {
	inter := runtime.NewIsolatedInterpreter(referenceFile, importRoots...)
	inter.ExternalDefinitions["os"] = extos.New(inter)
	inter.ExternalDefinitions["rx"] = extrx.New(inter)
	inter.ExternalDefinitions["docs"] = extdocs.New(inter)
	inter.ExternalDefinitions["fs"] = extfs.New(inter)
	return inter
}

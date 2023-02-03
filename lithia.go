package lithia

import (
	"context"
	"fmt"
	"strings"
	"time"

	extdocs "github.com/vknabel/lithia/external/docs"
	extfs "github.com/vknabel/lithia/external/fs"
	extos "github.com/vknabel/lithia/external/os"
	extrx "github.com/vknabel/lithia/external/rx"
	"github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/world"
)

func NewDefaultInterpreter(referenceFile string, importRoots ...string) (*runtime.Interpreter, context.Context) {
	ctx := context.Background()
	if timeout, ok := world.Current.Env.LookupEnv("LITHIA_TIMEOUT"); ok {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			fmt.Fprint(world.Current.Stderr, err)
			world.Current.Env.Exit(1)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, duration)

		go func() {
			<-time.After(duration)
			cancel()
		}()
	}
	inter := runtime.NewIsolatedInterpreter(ctx, referenceFile, importRoots...)

	externalDefinitions := "*"
	if defs, ok := world.Current.Env.LookupEnv("LITHIA_EXTERNAL_DEFINITIONS"); ok {
		externalDefinitions = defs
	}
	var whitelistExternals []string
	if externalDefinitions != "*" {
		whitelistExternals = strings.Split(externalDefinitions, ",")
	}

	if allowsWhitelistExternal(whitelistExternals, "os") {
		inter.ExternalDefinitions["os"] = extos.New(inter)
	}
	if allowsWhitelistExternal(whitelistExternals, "rx") {
		inter.ExternalDefinitions["rx"] = extrx.New(inter)
	}
	if allowsWhitelistExternal(whitelistExternals, "docs") {
		inter.ExternalDefinitions["docs"] = extdocs.New(inter)
	}
	if allowsWhitelistExternal(whitelistExternals, "fs") {
		inter.ExternalDefinitions["fs"] = extfs.New(inter)
	}
	return inter, ctx
}

func allowsWhitelistExternal(whitelistExternals []string, name string) bool {
	if whitelistExternals == nil {
		return true
	}
	for _, allowed := range whitelistExternals {
		if allowed == name {
			return true
		}
	}
	return false
}

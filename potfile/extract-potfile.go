package potfile

import (
	"fmt"
	"path"

	"github.com/vknabel/lithia"
	"github.com/vknabel/lithia/external/rx"
	"github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/world"
)

func ForReferenceFile(fileName string) (PotfileState, error) {
	inter := lithia.NewDefaultInterpreter(path.Dir(fileName))
	pkg := inter.Resolver.ResolvePackageForReferenceFile(fileName)
	if pkg.Manifest == nil {
		return PotfileState{}, fmt.Errorf("failed to find manifest for package")
	}
	return serializePotfile(pkg.Manifest.Path)
}

func serializePotfile(fileName string) (PotfileState, error) {
	scriptData, err := world.Current.FS.ReadFile(fileName)
	if err != nil {
		return PotfileState{}, err
	}
	inter := lithia.NewDefaultInterpreter(path.Dir(fileName))
	script := string(scriptData)
	_, err = inter.Interpret(fileName, script)
	if err != nil {
		return PotfileState{}, err
	}

	storeModule, ok := inter.Modules["pot"]
	if !ok {
		return PotfileState{}, fmt.Errorf("failed to find pot module")
	}
	lazyStore, ok := storeModule.Environment.GetExported("store")
	if !ok {
		return PotfileState{}, fmt.Errorf("failed to find store")
	}
	store, runErr := lazyStore.Evaluate()
	if err != nil {
		return PotfileState{}, runErr
	}
	storeTypeRef := store.RuntimeType()
	if storeTypeRef.Module != "apps" || storeTypeRef.Name != "Store" {
		return PotfileState{}, fmt.Errorf("failed to find store")
	}
	lazyStates, runErr := store.Lookup("states")
	if err != nil {
		return PotfileState{}, runErr
	}
	states, runErr := lazyStates.Evaluate()
	if runErr != nil {
		return PotfileState{}, runErr
	}
	statesVar, ok := states.(rx.RxVariable)
	if !ok {
		return PotfileState{}, fmt.Errorf("failed to find states")
	}
	latestState, runErr := statesVar.Current()
	if runErr != nil {
		return PotfileState{}, runErr
	}
	latestStateTypeRef := latestState.RuntimeType()
	if latestStateTypeRef.Module != "pot" || latestStateTypeRef.Name != "State" {
		return PotfileState{}, fmt.Errorf("failed to find state")
	}
	lazyCmds, runErr := latestState.Lookup("cmds")
	if runErr != nil {
		return PotfileState{}, runErr
	}
	anyCmdsDict, runErr := lazyCmds.Evaluate()
	if runErr != nil {
		return PotfileState{}, runErr
	}
	cmdsDict, ok := anyCmdsDict.(runtime.PreludeDict)
	if !ok {
		return PotfileState{}, fmt.Errorf("failed to find cmds")
	}
	cmdsMap := cmdsDict.ToMap()

	cmds := make(map[string]PotfileCmd)
	for i, lazyCmd := range cmdsMap {
		cmd, runErr := lazyCmd.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		cmdTypeRef := cmd.RuntimeType()
		if cmdTypeRef.Module != "pot.cmds" || cmdTypeRef.Name != "Command" {
			return PotfileState{}, fmt.Errorf("failed to find cmd")
		}
		lazyName, runErr := cmd.Lookup("name")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		name, runErr := lazyName.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		nameStr, ok := name.(runtime.PreludeString)
		if !ok {
			return PotfileState{}, fmt.Errorf("failed to find name")
		}
		lazySummary, runErr := cmd.Lookup("summary")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		summary, runErr := lazySummary.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		summaryStr, ok := summary.(runtime.PreludeString)
		if !ok {
			return PotfileState{}, fmt.Errorf("failed to find summary")
		}
		lazyFlags, runErr := cmd.Lookup("flags")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		flagsRuntimeDict, runErr := lazyFlags.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		flagsPreludeDict := flagsRuntimeDict.(runtime.PreludeDict)
		flagsLazyMap := flagsPreludeDict.ToMap()

		lazyBin, runErr := cmd.Lookup("bin")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		bin, runErr := lazyBin.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		binStr, ok := bin.(runtime.PreludeString)
		if !ok {
			return PotfileState{}, fmt.Errorf("failed to find bin")
		}
		lazyEnvs, runErr := cmd.Lookup("envs")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		lazyEnvsDict, runErr := lazyEnvs.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		envsDict, ok := lazyEnvsDict.(runtime.PreludeDict)
		if !ok {
			return PotfileState{}, fmt.Errorf("failed to find envs")
		}
		evalEnvsMap := envsDict.ToMap()
		envsMap := make(map[string]string, len(evalEnvsMap))

		for k, lazyV := range evalEnvsMap {
			v, runErr := lazyV.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			str, ok := v.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find envs")
			}
			envsMap[k] = string(str)
		}

		lazyArgs, runErr := cmd.Lookup("args")
		if runErr != nil {
			return PotfileState{}, runErr
		}
		argsList, runErr := lazyArgs.Evaluate()
		if runErr != nil {
			return PotfileState{}, runErr
		}
		argsSlice, err := listRuntimeValueToSlice(argsList)
		if err != nil {
			return PotfileState{}, err
		}
		args := make([]string, len(argsSlice))
		for k, v := range argsSlice {
			str, ok := v.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find args")
			}
			args[k] = string(str)
		}

		flagsMap := make(map[string]PotfileFlag)
		for j, lazyFlag := range flagsLazyMap {
			flag, runErr := lazyFlag.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			flagTypeRef := flag.RuntimeType()
			if flagTypeRef.Module != "pot.cmds" || flagTypeRef.Name != "Flag" {
				return PotfileState{}, fmt.Errorf("failed to find flag, got %s", flagTypeRef)
			}
			lazyName, runErr := flag.Lookup("name")
			if runErr != nil {
				return PotfileState{}, runErr
			}
			name, runErr := lazyName.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			nameStr, ok := name.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find name")
			}
			lazyShort, runErr := flag.Lookup("short")
			if runErr != nil {
				return PotfileState{}, runErr
			}
			short, runErr := lazyShort.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			shortStr, ok := short.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find short")
			}
			lazySummary, runErr := flag.Lookup("summary")
			if runErr != nil {
				return PotfileState{}, runErr
			}
			summary, runErr := lazySummary.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			summaryStr, ok := summary.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find summary")
			}
			lazyDefault, runErr := flag.Lookup("default")
			if runErr != nil {
				return PotfileState{}, runErr
			}
			defaultVal, runErr := lazyDefault.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			defaultStr, ok := defaultVal.(runtime.PreludeString)
			if !ok {
				return PotfileState{}, fmt.Errorf("failed to find default")
			}
			required, runErr := flag.Lookup("required")
			if runErr != nil {
				return PotfileState{}, runErr
			}
			requiredVal, runErr := required.Evaluate()
			if runErr != nil {
				return PotfileState{}, runErr
			}
			if requiredVal.RuntimeType().Name == "True" && requiredVal.RuntimeType().Module == "prelude" {
				flagsMap[j] = PotfileFlag{
					Name:         string(nameStr),
					Short:        string(shortStr),
					Summary:      string(summaryStr),
					DefaultValue: string(defaultStr),
					Required:     false,
				}
				continue
			}
			flagsMap[j] = PotfileFlag{
				Name:         string(nameStr),
				Short:        string(shortStr),
				Summary:      string(summaryStr),
				DefaultValue: string(defaultStr),
				Required:     false,
			}
		}
		cmds[i] = PotfileCmd{
			Name:    string(nameStr),
			Summary: string(summaryStr),
			Flags:   flagsMap,
			Bin:     string(binStr),
			Envs:    envsMap,
			Args:    args,
		}
	}
	return PotfileState{
		Cmds: cmds,
	}, nil
}

func listRuntimeValueToSlice(lv runtime.RuntimeValue) ([]runtime.RuntimeValue, error) {
	lvTypeRef := lv.RuntimeType()
	if lvTypeRef.Module != "prelude" {
		return nil, fmt.Errorf("expected list, got %s", lvTypeRef)
	}
	if lvTypeRef.Name == "Nil" {
		return []runtime.RuntimeValue{}, nil
	}
	if lvTypeRef.Name != "Cons" {
		return nil, fmt.Errorf("expected list, got %s", lvTypeRef)
	}
	lazyHead, runErr := lv.Lookup("head")
	if runErr != nil {
		return nil, *runErr
	}
	head, runErr := lazyHead.Evaluate()
	if runErr != nil {
		return nil, *runErr
	}
	lazyTail, runErr := lv.Lookup("tail")
	if runErr != nil {
		return nil, *runErr
	}
	tail, runErr := lazyTail.Evaluate()
	if runErr != nil {
		return nil, *runErr
	}
	sliceTail, err := listRuntimeValueToSlice(tail)
	if err != nil {
		return nil, err
	}
	return append([]runtime.RuntimeValue{head}, sliceTail...), nil
}

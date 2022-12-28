package runtime

import (
	"fmt"
	"sort"
	"strings"
)

var _ RuntimeValue = PreludeDict{}
var PreludeDictTypeRef = MakeRuntimeTypeRef("Dict", "prelude")

type PreludeDict struct {
	dict    map[PreludeString]Evaluatable
	context *InterpreterContext
}

func MakePreludeDict(context *InterpreterContext, dict map[PreludeString]Evaluatable) PreludeDict {
	copy := map[PreludeString]Evaluatable{}
	for k, v := range dict {
		copy[k] = v
	}
	return PreludeDict{
		dict:    copy,
		context: context,
	}
}

func (d PreludeDict) EagerEvaluate() *RuntimeError {
	for _, m := range d.dict {
		value, err := m.Evaluate()
		if err != nil {
			return err
		}
		eagerEvaluatable, ok := value.(EagerEvaluatableRuntimeValue)
		if !ok {
			continue
		}
		err = eagerEvaluatable.EagerEvaluate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (PreludeDict) RuntimeType() RuntimeTypeRef {
	return PreludeDictTypeRef
}

func (rv PreludeDict) String() string {
	if len(rv.dict) == 0 {
		return "[:]"
	}
	entries := []string{}
	for k, v := range rv.dict {
		value, err := v.Evaluate()
		if err != nil {
			entries = append(entries, fmt.Sprintf("%s: %s", k, err.Error()))
		} else {
			entries = append(entries, fmt.Sprintf("%s: %s", k, value))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(entries, ", "))
}

func (rv PreludeDict) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "length":
		return NewConstantRuntimeValue(PreludeInt(len(rv.dict))), nil
	case "get":
		return NewConstantRuntimeValue(MakeAnonymousFunction(
			"get",
			[]string{"key"},
			func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
				key, err := args[0].Evaluate()
				if err != nil {
					return nil, err
				}
				if stringKey, ok := key.(PreludeString); ok {
					if value, ok := rv.dict[stringKey]; ok {
						return rv.context.environment.MakeSome(value)
					} else {
						return rv.context.environment.MakeNone()
					}
				} else {
					return nil, NewRuntimeErrorf("dict key must be a string, got %s", key)
				}
			})), nil
	case "set":
		return NewConstantRuntimeValue(MakeAnonymousFunction(
			"set",
			[]string{"key", "value"},
			func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
				key, err := args[0].Evaluate()
				if err != nil {
					return nil, err
				}

				if stringKey, ok := key.(PreludeString); ok {
					copy := rv.copy()
					copy.dict[stringKey] = args[1]
					return copy, nil
				} else {
					return nil, NewRuntimeErrorf("dict key must be a string, got %s", key)
				}
			})), nil
	case "delete":
		return NewConstantRuntimeValue(MakeAnonymousFunction(
			"delete",
			[]string{"key"},
			func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
				key, err := args[0].Evaluate()
				if err != nil {
					return nil, err
				}

				if stringKey, ok := key.(PreludeString); ok {
					copy := rv.copy()
					delete(copy.dict, stringKey)
					return copy, nil
				} else {
					return nil, NewRuntimeErrorf("dict key must be a string, got %s", key)
				}
			})), nil
	case "entries":
		pairs := make([]RuntimeValue, 0, len(rv.dict))
		for _, k := range rv.orderedKeys() {
			v := rv.dict[k]
			pair, err := rv.context.environment.MakePair(NewConstantRuntimeValue(k), v)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, pair)
		}
		dataList, err := rv.context.environment.MakeEagerList(pairs)
		return NewConstantRuntimeValue(dataList), err
	case "keys":
		keys := make([]RuntimeValue, 0, len(rv.dict))
		for _, k := range rv.orderedKeys() {
			keys = append(keys, k)
		}
		dataList, err := rv.context.environment.MakeEagerList(keys)
		return NewConstantRuntimeValue(dataList), err
	case "values":
		values := make([]Evaluatable, 0, len(rv.dict))
		for _, k := range rv.orderedKeys() {
			values = append(values, rv.dict[k])
		}
		dataList, err := rv.context.environment.MakeList(values)
		return NewConstantRuntimeValue(dataList), err
	default:
		return nil, NewRuntimeErrorf("no such member: %s for %s", member, rv.RuntimeType().String())
	}
}

func (rv PreludeDict) ToMap() map[string]Evaluatable {
	result := make(map[string]Evaluatable)
	for k, v := range rv.dict {
		result[string(k)] = v
	}
	return result
}

func (rv PreludeDict) copy() PreludeDict {
	return MakePreludeDict(rv.context, rv.dict)
}

func (rv PreludeDict) orderedKeys() []PreludeString {
	keys := make([]PreludeString, 0, len(rv.dict))
	for k := range rv.dict {
		keys = append(keys, k)
	}
	sort.Sort(preludeStringSlice(keys))
	return keys
}

type preludeStringSlice []PreludeString

func (x preludeStringSlice) Len() int           { return len(x) }
func (x preludeStringSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x preludeStringSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

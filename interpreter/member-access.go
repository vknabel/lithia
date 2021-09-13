package interpreter

import "fmt"

func (dataValue DataRuntimeValue) Lookup(name string) (*LazyRuntimeValue, error) {
	if lazyValue, ok := dataValue.members[name]; ok {
		return lazyValue, nil
	} else {
		return nil, fmt.Errorf("%s is not a member of %s", name, dataValue)
	}
}

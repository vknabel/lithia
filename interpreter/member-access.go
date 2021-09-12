package interpreter

import "fmt"

type MemberAccessable interface {
	Lookup(string) (RuntimeValue, error)
}

func (dataValue DataRuntimeValue) Lookup(name string) (RuntimeValue, error) {
	if lazyValue, ok := dataValue.members[name]; ok {
		value, err := lazyValue.Evaluate()
		if err != nil {
			return nil, err
		}
		return value, nil
	} else {
		return nil, fmt.Errorf("%s is not a member of %s", name, dataValue)
	}
}

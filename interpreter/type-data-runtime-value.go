package interpreter

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = DataRuntimeValue{}

type DataRuntimeValue struct {
	typeValue *DataDeclRuntimeValue
	members   map[string]*LazyRuntimeValue
}

func (d DataRuntimeValue) RuntimeType() RuntimeType {
	var typeValue RuntimeValue = *d.typeValue
	return RuntimeType{
		name:      d.typeValue.name,
		typeValue: &typeValue,
	}
}

func (d DataRuntimeValue) String() string {
	params := make([]string, 0)
	for _, arg := range d.members {
		value, err := arg.Evaluate()
		if err != nil {
			params = append(params, err.Error())
		} else {
			params = append(params, fmt.Sprint(value))
		}
	}

	return fmt.Sprintf("%s %s", d.typeValue.name, strings.Join(params, ", "))
}

package interpreter

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = DataDeclRuntimeValue{}
var _ Callable = DataDeclRuntimeValue{}

type DataDeclRuntimeValue struct {
	name   string
	fields []DataDeclField
}

type DataDeclField struct {
	name   string
	params []string
}

func (d DataDeclRuntimeValue) String() string {
	fieldNames := make([]string, len(d.fields))
	for i, field := range d.fields {
		fieldNames[i] = field.name
		if len(field.params) > 0 {
			fieldNames[i] += " " + strings.Join(field.params, ", ")
		}
	}
	return fmt.Sprintf("data %s { %s }", d.name, strings.Join(fieldNames, "; "))
}

func (d DataDeclRuntimeValue) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (f DataDeclRuntimeValue) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("data type %s has no member %s", f, member)
}

func (dataDecl DataDeclRuntimeValue) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	if len(arguments) < len(dataDecl.fields) {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
			return CurriedCallable{
				actual:         dataDecl,
				args:           arguments,
				remainingArity: len(dataDecl.fields) - len(arguments),
			}, nil
		})
		return lazy.Evaluate()
	} else if len(arguments) == len(dataDecl.fields) {
		members := make(map[string]*LazyRuntimeValue, len(dataDecl.fields))
		for i, field := range dataDecl.fields {
			arg := arguments[i]
			members[field.name] = arg
		}
		instance := DataRuntimeValue{
			typeValue: &dataDecl,
			members:   members,
		}
		return instance, nil
	} else {
		// error
		return nil, fmt.Errorf("too many arguments")
	}
}

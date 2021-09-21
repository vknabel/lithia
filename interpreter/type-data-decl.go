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
	docs   DocString
}

func NewDataDecl(name string, fields []DataDeclField, docs DocString) DataDeclRuntimeValue {
	return DataDeclRuntimeValue{
		name:   name,
		fields: fields,
		docs:   docs,
	}
}

type DataDeclField struct {
	name   string
	params []string
	docs   DocString
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

func (f DataDeclRuntimeValue) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("data type %s has no member %s", f, member)
}

func (dataDecl DataDeclRuntimeValue) Call(arguments []Evaluatable) (RuntimeValue, error) {
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
		members := make(map[string]Evaluatable, len(dataDecl.fields))
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

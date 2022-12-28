package runtime

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = DataRuntimeValue{}

type DataRuntimeValue struct {
	TypeDecl *PreludeDataDecl
	Members  map[string]Evaluatable
}

func MakeDataRuntimeValueMemberwise(decl *PreludeDataDecl, members map[string]Evaluatable) (DataRuntimeValue, *RuntimeError) {
	if len(members) != len(decl.Decl.Fields) {
		return DataRuntimeValue{}, NewRuntimeErrorf("wrong number of members")
	}

	copiedMembers := make(map[string]Evaluatable)
	for _, field := range decl.Decl.Fields {
		if value, ok := members[string(field.Name)]; ok {
			copiedMembers[string(field.Name)] = value
		} else {
			return DataRuntimeValue{}, NewRuntimeErrorf("missing %s", field.Name)
		}
	}
	return DataRuntimeValue{
		TypeDecl: decl,
		Members:  copiedMembers,
	}, nil
}

func (d DataRuntimeValue) String() string {
	params := make([]string, 0)
	for _, field := range d.TypeDecl.Decl.Fields {
		arg := d.Members[string(field.Name)]
		value, err := arg.Evaluate()
		if err != nil {
			params = append(params, err.Error())
		} else {
			params = append(params, fmt.Sprint(value))
		}
	}

	return fmt.Sprintf("(%s %s)", d.TypeDecl.Decl.Name, strings.Join(params, ", "))
}

func (d DataRuntimeValue) EagerEvaluate() *RuntimeError {
	for _, m := range d.Members {
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

func (d DataRuntimeValue) Lookup(name string) (Evaluatable, *RuntimeError) {
	if value, ok := d.Members[name]; ok {
		return value, nil
	} else {
		return nil, NewRuntimeErrorf("cannot read property %s of %s", name, d)
	}
}

func (d DataRuntimeValue) RuntimeType() RuntimeTypeRef {
	return RuntimeTypeRef{
		Name:   d.TypeDecl.Decl.Name,
		Module: d.TypeDecl.Decl.MetaInfo.ModuleName,
	}
	// d.TypeDecl.Decl.DeclName()
	// panic("TODO: there are not only global dependencies!")
}

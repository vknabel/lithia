package interpreter

import "fmt"

type Environment struct {
	Parent     *Environment
	Scope      map[string]*LazyRuntimeValue
	Unexported map[string]*LazyRuntimeValue
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		Parent:     parent,
		Scope:      make(map[string]*LazyRuntimeValue),
		Unexported: make(map[string]*LazyRuntimeValue),
	}
}

func (env *Environment) Private() *Environment {
	return &Environment{
		Parent:     env.Parent,
		Scope:      env.Scope,
		Unexported: make(map[string]*LazyRuntimeValue),
	}
}

func (env *Environment) Declare(name string, value *LazyRuntimeValue) error {
	if env.DirectlyDefines(name) {
		return fmt.Errorf("variable %s already declared", name)
	}
	env.Scope[name] = value
	return nil
}

func (env *Environment) DeclareUnexported(name string, value *LazyRuntimeValue) error {
	if env.DirectlyDefines(name) {
		return fmt.Errorf("variable %s already declared", name)
	}
	env.Unexported[name] = value
	return nil
}

func (env *Environment) Get(name string) (*LazyRuntimeValue, bool) {
	if value, ok := env.Scope[name]; ok {
		return value, true
	}
	if value, ok := env.Unexported[name]; ok {
		return value, true
	}

	if env.Parent != nil {
		return env.Parent.Get(name)
	}

	return nil, false
}

func (env *Environment) DirectlyDefines(name string) bool {
	if _, ok := env.Scope[name]; ok {
		return true
	}
	if _, ok := env.Unexported[name]; ok {
		return true
	}
	return false
}

func (env *Environment) Contains(name string) bool {
	if env.DirectlyDefines(name) {
		return true
	}
	if env.Parent != nil {
		return env.Parent.Contains(name)
	}
	return false
}

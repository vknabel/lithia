package interpreter

import "fmt"

type Environment struct {
	Parent *Environment
	Scope  map[string]*LazyRuntimeValue
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		Parent: parent,
		Scope:  make(map[string]*LazyRuntimeValue),
	}
}

func (env *Environment) Declare(name string, value *LazyRuntimeValue) error {
	if _, ok := env.Scope[name]; ok {
		return fmt.Errorf("variable %s already declared", name)
	}
	env.Scope[name] = value
	return nil
}

func (env *Environment) Get(name string) (*LazyRuntimeValue, bool) {
	if value, ok := env.Scope[name]; ok {
		return value, true
	}

	if env.Parent != nil {
		return env.Parent.Get(name)
	}

	return nil, false
}

func (env *Environment) Contains(name string) bool {
	if _, ok := env.Scope[name]; ok {
		return true
	}

	if env.Parent != nil {
		return env.Parent.Contains(name)
	}

	return false
}

package ast

type ModuleName string

type ContextModule struct {
	Name ModuleName

	Files []*SourceFile
}

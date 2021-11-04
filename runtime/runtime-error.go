package runtime

import "github.com/vknabel/go-lithia/ast"

type RuntimeError struct {
	Cause      error
	StackTrace []ast.Source
}

func NewRuntimeError(err error, source ast.Source) *RuntimeError {
	return &RuntimeError{
		Cause:      err,
		StackTrace: []ast.Source{source},
	}
}

func (r RuntimeError) Cascade(source ast.Source) *RuntimeError {
	return &RuntimeError{
		Cause:      r.Cause,
		StackTrace: append(r.StackTrace, source),
	}
}

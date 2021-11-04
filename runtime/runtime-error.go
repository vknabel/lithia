package runtime

import "github.com/vknabel/go-lithia/ast"

type RuntimeError struct {
	Cause      error
	StackTrace []ast.Source
}

func NewRuntimeError(err error) *RuntimeError {
	return &RuntimeError{
		Cause:      err,
		StackTrace: []ast.Source{},
	}
}

func (r *RuntimeError) Cascade(source ast.Source) *RuntimeError {
	if r == nil {
		return nil
	}
	return &RuntimeError{
		Cause:      r.Cause,
		StackTrace: append(r.StackTrace, source),
	}
}

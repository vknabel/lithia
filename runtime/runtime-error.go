package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

type RuntimeError struct {
	Cause      error
	StackTrace []ast.Source
}

func (e RuntimeError) Error() string {
	// TODO: Stacktrace
	return fmt.Sprintf("runtime error: %s", e.Cause)
}

func NewRuntimeError(err error) *RuntimeError {
	if runtimeErr, ok := err.(*RuntimeError); ok {
		return runtimeErr
	}
	return &RuntimeError{
		Cause:      err,
		StackTrace: []ast.Source{},
	}
}

func NewRuntimeErrorf(format string, args ...interface{}) *RuntimeError {
	return NewRuntimeError(fmt.Errorf(format, args...))
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

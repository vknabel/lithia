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
	stackTrace := ""
	for _, source := range e.StackTrace {
		stackTrace += fmt.Sprintf(
			"\t%s:%d:%d %s\n",
			source.FileName,
			source.Start.Line+1,
			source.Start.Column+1,
			source.ModuleName,
		)
	}
	return fmt.Sprintf("runtime error: %s\n%s", e.Cause, stackTrace)
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
	if len(r.StackTrace) > 0 && r.StackTrace[len(r.StackTrace)-1] == source {
		return r
	}
	return &RuntimeError{
		Cause:      r.Cause,
		StackTrace: append(r.StackTrace, source),
	}
}

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
		stackTrace += stackTraceEntryString(source)
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
	if len(r.StackTrace) == 0 {
		return &RuntimeError{
			Cause:      r.Cause,
			StackTrace: append(r.StackTrace, source),
		}
	}
	last := r.StackTrace[len(r.StackTrace)-1]
	if stackTraceEntryString(last) == stackTraceEntryString(source) {
		return r
	} else {
		return &RuntimeError{
			Cause:      r.Cause,
			StackTrace: append(r.StackTrace, source),
		}
	}
}

func (r *RuntimeError) CascadeDecl(decl ast.Decl) *RuntimeError {
	if decl.Meta().Source == nil {
		return r
	} else {
		return r.Cascade(*decl.Meta().Source)
	}
}

func stackTraceEntryString(source ast.Source) string {
	return fmt.Sprintf(
		"\t%s:%d:%d %s\n",
		source.FileName,
		source.Start.Line+1,
		source.Start.Column+1,
		source.ModuleName,
	)
}

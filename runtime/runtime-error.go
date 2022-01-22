package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vknabel/go-lithia/ast"
)

type RuntimeError struct {
	topic      string
	cause      error
	stackTrace []stackEntry
}

type stackEntry struct {
	// contextPath []string
	source ast.Source

	// optional
	decl ast.Decl
	// optional
	expr ast.Expr
}

func (e RuntimeError) Error() string {
	stackTrace := ""
	for _, source := range e.stackTrace {
		stackTrace += stackTraceEntryString(source)
	}
	return fmt.Sprintf("%s: %s\n%s", e.topic, e.cause, stackTrace)
}

func NewRuntimeError(err error) *RuntimeError {
	if runtimeErr, ok := err.(*RuntimeError); ok {
		return runtimeErr
	}
	return &RuntimeError{
		topic:      "runtime error",
		cause:      err,
		stackTrace: []stackEntry{},
	}
}

func NewRuntimeErrorf(format string, args ...interface{}) *RuntimeError {
	return NewRuntimeError(fmt.Errorf(format, args...))
}

func (r *RuntimeError) cascadeEntry(entry stackEntry) *RuntimeError {
	if r == nil {
		return nil
	}
	if len(r.stackTrace) == 0 {
		return &RuntimeError{
			topic:      r.topic,
			cause:      r.cause,
			stackTrace: append(r.stackTrace, entry),
		}
	}
	last := r.stackTrace[len(r.stackTrace)-1]
	if stackTraceEntryString(last) == stackTraceEntryString(entry) {
		return r
	} else {
		return &RuntimeError{
			topic:      r.topic,
			cause:      r.cause,
			stackTrace: append(r.stackTrace, entry),
		}
	}
}

func (r *RuntimeError) CascadeDecl(decl ast.Decl) *RuntimeError {
	if decl.Meta().Source == nil {
		return r
	} else {
		return r.cascadeEntry(stackEntry{
			source: *decl.Meta().Source,
			decl:   decl,
		})
	}
}

func (r *RuntimeError) CascadeExpr(expr ast.Expr) *RuntimeError {
	if expr.Meta().Source == nil {
		return r
	} else {
		return r.cascadeEntry(stackEntry{
			source: *expr.Meta().Source,
			expr:   expr,
		})
	}
}

func stackTraceEntryString(entry stackEntry) string {
	var name string
	if entry.decl != nil {
		name = string(entry.decl.DeclName())
	}

	fileName := entry.source.FileName
	if dir, err := os.Getwd(); err == nil {
		rel, err := filepath.Rel(dir, entry.source.FileName)
		if err == nil {
			fileName = "." + string(os.PathSeparator) + rel
		}
	}
	source := entry.source
	return fmt.Sprintf(
		"\t%s:%d:%d %s\n",
		fileName,
		source.Start.Line+1,
		source.Start.Column+1,
		name,
	)
}

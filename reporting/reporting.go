package reporting

import (
	"fmt"

	"github.com/vknabel/lithia/world"
)

func ReportErrorOrPanic(err error) {
	fmt.Fprintln(world.Current.Stderr, err)
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintln(world.Current.Stderr, "[line"+fmt.Sprint(line)+"] Error"+where+": "+message)
}

type LocatableError interface {
	error
	SourceLocation() string
}

package reporting

import (
	"fmt"
	"os"
)

func ReportErrorOrPanic(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line"+fmt.Sprint(line)+"] Error"+where+": "+message)
}

type LocatableError interface {
	error
	SourceLocation() string
}

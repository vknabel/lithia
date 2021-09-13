package reporting

import (
	"fmt"
	"os"
)

func ReportErrorOrPanic(err error) {
	if locatable, ok := err.(LocatableError); ok {
		fmt.Fprintln(os.Stderr, locatable)
	} else {
		report(0, "", err.Error())
	}
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line"+fmt.Sprint(line)+"] Error"+where+": "+message)
}

type LocatableError interface {
	SourceLocation() string
}

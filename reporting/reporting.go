package reporting

import (
	"fmt"
	"os"
)

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line"+fmt.Sprint(line)+"] Error"+where+": "+message)
}

package interpreter

import (
	"strings"
)

type DocString string

func MakeDocString(comments []string) DocString {
	docsStrings := make([]string, 0)
	for _, comment := range comments {
		if strings.HasPrefix(comment, "///") {
			docsStrings = append(docsStrings, strings.TrimSpace(strings.TrimPrefix(comment, "///")))
		} else if strings.HasPrefix(comment, "/**") {
			trimmed := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "/**"), "*/"))
			lines := strings.Split(trimmed, "\n")
			for _, line := range lines {
				docsStrings = append(docsStrings, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*")))
			}
		}
	}
	return DocString(strings.Join(docsStrings, "\n"))
}

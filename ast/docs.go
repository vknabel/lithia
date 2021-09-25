package ast

import "strings"

type Docs struct {
	// not only text, but also parsing @type and @param, ...
	// or more general?
	Content string
}

func MakeDocs(comments []string) *Docs {
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
	return &Docs{
		Content: strings.Join(docsStrings, "\n"),
	}
}

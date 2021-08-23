package scanner

import "fmt"

type Token struct {
	Type                 TokenType
	Lexeme               string
	Literal              interface{}
	Line                 int
	Column               int
	Offset               int
	PrecedingAnnotations []TokenAnnotation
}

type AnnotationType int

const (
	ANNOTATION_WHITESPACE = iota
	ANNOTATION_LINE_COMMENT
	ANNOTATION_BLOCK_COMMENT
)

type TokenAnnotation struct {
	Lexeme string
	Type   AnnotationType
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s  %s", t.Type.String(), t.Lexeme, t.Literal)
}

package scanner

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
	Offset  int
	// TODO: later add surrounding whitespace and comments
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s  %s", t.Type.String(), t.Lexeme, t.Literal)
}

package scanner

import "fmt"

type TokenType int

const (
	// single character tokens
	COMMA                = iota // ,
	LEFT_PAREN                  // (
	RIGHT_PAREN                 // )
	LEFT_BRACE                  // {
	RIGHT_BRACE                 // }
	LEFT_BRACKET                // [
	RIGHT_BRACKET               // ]
	SEMICOLON_OR_NEWLINE        // ;
	COLON                       // :
	DOT                         // .
	MINUS                       // -
	PLUS                        // +
	STAR                        // *
	SLASH                       // /
	PERCENT                     // %

	// one or two characters
	BANG          // !
	BANG_EQUAL    // !=
	EQUAL         // ==
	EQUAL_EQUAL   // ===
	ARROW         // =>
	GREATER       // >
	GREATER_EQUAL // >=
	LESS          // <
	LESS_EQUAL    // <=

	// literals
	IDENTIFIER // identifier
	STRING     // "string"
	NUMBER     // 3.14159

	// keywords
	DATA   // data Type { }
	UNION  // union Type { }
	ENUM   // enum Type { }
	LET    // let var = expr
	FUNC   // func tion { _ => }
	IMPORT // import package

	// special tokens
	EOF
	ILLEGAL
	LINE_COMMENT
	BLOCK_COMMENT
)

func (t TokenType) String() string {
	switch t {
	case COMMA:
		return "COMMA"
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case SEMICOLON_OR_NEWLINE:
		return "SEMICOLON_OR_NEWLINE"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case PERCENT:
		return "PERCENT"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case ARROW:
		return "ARROW"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case DATA:
		return "DATA"
	case UNION:
		return "UNION"
	case ENUM:
		return "ENUM"
	case LET:
		return "LET"
	case FUNC:
		return "FUNC"
	case IMPORT:
		return "IMPORT"
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	}
	return fmt.Sprintf("TOKEN#%d", t)
}

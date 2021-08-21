package scanner

import (
	"bufio"
	"fmt"
	"io"
)

type Scanner struct {
	reader *bufio.Reader

	currentLexeme string
	offset        int
	line          int
	column        int
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
		offset: 0,
		line:   1,
		column: 0,
	}
}

func (scanner *Scanner) ScanTokens() ([]Token, []error) {
	tokens := make([]Token, 0)
	errors := make([]error, 0)
	for {
		// scanner.start = scanner.current

		tok, err := scanner.scanToken()

		if err != nil {
			return tokens, []error{err}
		}
		if tok.Type == EOF {
			return tokens, errors
		} else if tok.Type == BLOCK_COMMENT || tok.Type == LINE_COMMENT {
			continue
		} else if tok.Type == ILLEGAL {
			errors = append(errors, fmt.Errorf("illegal character '%s'", tok.Lexeme))
		} else {
			tokens = append(tokens, tok)
		}
	}
}

func (scanner *Scanner) advance() (rune, int, error) {
	rune, size, err := scanner.reader.ReadRune()
	scanner.offset += size
	if rune == '\n' {
		scanner.advanceLine()
	}
	return rune, size, err
}

func (scanner *Scanner) peek() (rune, int, error) {
	rune, size, err := scanner.reader.ReadRune()
	if err != nil {
		return rune, size, err
	}
	err = scanner.reader.UnreadRune()
	return rune, size, err
}

func (scanner *Scanner) backtrack(rune rune) {
	scanner.reader.UnreadRune()
	scanner.offset -= 1
	if rune == '\n' {
		scanner.line -= 1
		scanner.column = 0 // TODO: Bug around newlines
	} else {
		scanner.column -= 1
	}
}

func (scanner *Scanner) match(char rune) bool {
	r, _, err := scanner.advance()
	if err != nil {
		return false
	}
	if r == char {
		return true
	} else {
		scanner.backtrack(r)
		return false
	}
}

func (scanner *Scanner) addToken(tokenType TokenType) Token {
	return scanner.addTokenLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Lexeme:  scanner.currentLexeme,
		Line:    scanner.line,
		Column:  scanner.column,
		Offset:  scanner.offset,
	}
}

func (scanner *Scanner) advanceLine() {
	scanner.line++
	scanner.column = 0
}

package scanner

import (
	"bufio"
	"fmt"
	"io"
)

type Scanner struct {
	reader *bufio.Reader
	tokens []Token

	annotations   []TokenAnnotation
	currentLexeme string
	offset        int
	line          int
	column        int
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader:      bufio.NewReader(reader),
		tokens:      make([]Token, 0),
		offset:      0,
		line:        1,
		column:      0,
		annotations: make([]TokenAnnotation, 0),
	}
}

func (scanner *Scanner) ScanTokens() ([]Token, []error) {
	tokens := make([]Token, 0)
	errors := make([]error, 0)
	currentLen := 0
	for {
		err := scanner.scanToken()
		if err != nil {
			return scanner.tokens, []error{err}
		}

		if currentLen == len(scanner.tokens) {
			continue
		}
		currentLen = len(scanner.tokens)
		token := scanner.tokens[len(scanner.tokens)-1]

		if token.Type == EOF {
			return tokens, errors
		} else if token.Type == ILLEGAL {
			errors = append(errors, fmt.Errorf("illegal character '%s'", token.Lexeme))
		} else {
			tokens = append(tokens, token)
		}
	}
}

func (scanner *Scanner) advance() (rune, int, error) {
	rune, size, err := scanner.reader.ReadRune()
	scanner.currentLexeme += string(rune)
	scanner.offset += 1
	if rune == '\n' {
		scanner.advanceLine()
	} else {
		scanner.column += 1
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
	if err := scanner.reader.UnreadRune(); err != nil {
		panic(err)
	}
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

func (scanner *Scanner) addToken(tokenType TokenType) {
	scanner.addTokenLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	lexeme := scanner.currentLexeme
	scanner.currentLexeme = ""
	annotations := scanner.annotations
	scanner.annotations = make([]TokenAnnotation, 0)

	scanner.tokens = append(scanner.tokens, Token{
		Type:                 tokenType,
		Literal:              literal,
		Lexeme:               lexeme,
		Line:                 scanner.line,
		Column:               scanner.column,
		Offset:               scanner.offset,
		PrecedingAnnotations: annotations,
	})
}

func (scanner *Scanner) addAnnotation(annotationType AnnotationType) {
	lexeme := scanner.currentLexeme
	scanner.currentLexeme = ""
	scanner.annotations = append(scanner.annotations, TokenAnnotation{
		Type:   annotationType,
		Lexeme: lexeme,
	})
}

func (scanner *Scanner) advanceLine() {
	scanner.line++
	scanner.column = 0
}

package scanner

import (
	"fmt"
	"io"
	"strconv"
)

// TODO: remove token return and add comments to scanner and tokens instead
func (scanner *Scanner) scanToken() error {
	rune, _, err := scanner.advance()

	if err == io.EOF {
		scanner.addToken(EOF)
		return nil
	}
	if err != nil {
		scanner.addTokenLiteral(ILLEGAL, err)
	}

	switch rune {
	case ',':
		scanner.addToken(COMMA)
		return nil
	case '(':
		scanner.addToken(LEFT_PAREN)
		return nil
	case ')':
		scanner.addToken(RIGHT_PAREN)
		return nil
	case '{':
		scanner.addToken(LEFT_BRACE)
		return nil
	case '}':
		scanner.addToken(RIGHT_BRACE)
		return nil
	case '[':
		scanner.addToken(LEFT_BRACKET)
		return nil
	case ']':
		scanner.addToken(RIGHT_BRACKET)
		return nil
	case ';':
		scanner.addToken(SEMICOLON)
		return nil
	case '\n':
		scanner.addToken(NEWLINE)
		return nil
	case ':':
		scanner.addToken(COLON)
		return nil
	case '.':
		scanner.addToken(DOT)
		return nil
	case '-':
		scanner.addToken(MINUS)
		return nil
	case '+':
		scanner.addToken(PLUS)
		return nil
	case '*':
		scanner.addToken(STAR)
		return nil
	case '%':
		scanner.addToken(PERCENT)
		return nil
	case '!':
		if scanner.match('=') {
			scanner.addToken(BANG_EQUAL)
			return nil
		} else {
			scanner.addToken(BANG)
			return nil
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(EQUAL_EQUAL)
			return nil
		} else if scanner.match('>') {
			scanner.addToken(ARROW)
			return nil
		} else {
			scanner.addToken(EQUAL)
			return nil
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(LESS_EQUAL)
			return nil
		} else {
			scanner.addToken(LESS)
			return nil
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(GREATER_EQUAL)
			return nil
		} else {
			scanner.addToken(GREATER)
			return nil
		}
	case '/':
		if scanner.match('/') {
			return scanner.scanLineComment()
		} else if scanner.match('*') {
			return scanner.scanMultilineComment()
		} else {
			scanner.addToken(SLASH)
			return nil
		}
	case ' ', '\t':
		scanner.addAnnotation(ANNOTATION_WHITESPACE)
		return nil
	case '"':
		return scanner.scanString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return scanner.scanNumber()
	default:
		if isLetter(rune) {
			return scanner.scanIdentifier()
		} else {
			scanner.addToken(ILLEGAL)
			return nil
		}
	}
}

func (scanner *Scanner) scanLineComment() error {
	for {
		rune, _, err := scanner.advance()
		if err == io.EOF {
			scanner.addAnnotation(ANNOTATION_LINE_COMMENT)
			return nil
		} else if err != nil {
			return err
		} else if rune == '\n' {
			scanner.addAnnotation(ANNOTATION_LINE_COMMENT)
			scanner.backtrack(rune)
			return nil
		}
	}
}

func (scanner *Scanner) scanMultilineComment() error {
	for {
		rune, _, err := scanner.advance()
		if err == io.EOF {
			scanner.addAnnotation(ANNOTATION_BLOCK_COMMENT)
			return nil
		} else if err != nil {
			return err
		} else if rune == '*' && scanner.match('/') {
			scanner.addAnnotation(ANNOTATION_BLOCK_COMMENT)
			return nil
		}
	}
}

func (scanner *Scanner) scanString() error {
	for {
		char, _, err := scanner.peek()
		if err == io.EOF {
			return fmt.Errorf("unexpected EOF")
		} else if err != nil {
			return err
		} else if char == '"' {
			scanner.match(char)
			value, err := strconv.Unquote(scanner.currentLexeme)
			if err != nil {
				return err
			}
			scanner.addTokenLiteral(STRING, value)
			return nil
		}
	}
}

func (scanner *Scanner) scanNumber() error {
	for {
		rune, _, err := scanner.advance()
		if err == io.EOF {
			value, err := strconv.Atoi(scanner.currentLexeme)
			if err != nil {
				return err
			}
			scanner.addTokenLiteral(INT, value)
			return nil
		} else if err != nil {
			return err
		} else if rune == '.' {
			// INT have no members!
			return scanner.scanFloat()
		} else if !isDigit(rune) {
			value, err := strconv.Atoi(scanner.currentLexeme)
			if err != nil {
				return err
			}
			scanner.addTokenLiteral(INT, value)
			return nil
		}
	}
}

func (scanner *Scanner) scanFloat() error {
	for {
		rune, _, err := scanner.peek()
		if err == io.EOF {
			value, err := strconv.ParseFloat(scanner.currentLexeme, 64)
			if err != nil {
				return err
			}
			scanner.addTokenLiteral(FLOAT, value)
			return nil
		} else if err != nil {
			return err
		} else if !isDigit(rune) {
			value, err := strconv.ParseFloat(scanner.currentLexeme, 64)
			if err != nil {
				return err
			}
			scanner.addTokenLiteral(FLOAT, value)
			return nil
		} else {
			scanner.advance()
		}
	}
}

func (scanner *Scanner) scanIdentifier() error {
	for {
		rune, _, err := scanner.peek()
		if err == io.EOF {
			scanner.addIdentifierTokenLiteral()
			return nil
		} else if err != nil {
			return err
		} else if !isAlphaNumeric(rune) {
			scanner.addIdentifierTokenLiteral()
			return nil
		} else {
			_, _, err = scanner.advance()
			if err != nil {
				return err
			}
		}
	}
}

func (scanner *Scanner) addIdentifierTokenLiteral() {
	if token, ok := reservedKeywords[scanner.currentLexeme]; ok {
		scanner.addTokenLiteral(token, scanner.currentLexeme)
	} else {
		scanner.addTokenLiteral(IDENTIFIER, scanner.currentLexeme)
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isDigit(r) || isLetter(r)
}

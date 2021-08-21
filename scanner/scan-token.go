package scanner

import "io"

// TODO: remove token return and add comments to scanner and tokens instead
func (scanner *Scanner) scanToken() (Token, error) {
	rune, _, err := scanner.advance()
	scanner.currentLexeme = string(rune)

	if err == io.EOF {
		return scanner.addToken(EOF), nil
	}
	if err != nil {
		scanner.addTokenLiteral(ILLEGAL, err)
	}

	switch rune {
	case ',':
		return scanner.addToken(COMMA), nil
	case '(':
		return scanner.addToken(LEFT_PAREN), nil
	case ')':
		return scanner.addToken(RIGHT_PAREN), nil
	case '{':
		return scanner.addToken(LEFT_BRACE), nil
	case '}':
		return scanner.addToken(RIGHT_BRACE), nil
	case '[':
		return scanner.addToken(LEFT_BRACKET), nil
	case ']':
		return scanner.addToken(RIGHT_BRACKET), nil
	case ';', '\n':
		return scanner.addToken(SEMICOLON_OR_NEWLINE), nil
	case ':':
		return scanner.addToken(COLON), nil
	case '.':
		return scanner.addToken(DOT), nil
	case '-':
		return scanner.addToken(MINUS), nil
	case '+':
		return scanner.addToken(PLUS), nil
	case '*':
		return scanner.addToken(STAR), nil
	case '%':
		return scanner.addToken(PERCENT), nil
	case '!':
		if scanner.match('=') {
			return scanner.addToken(BANG_EQUAL), nil
		} else {
			return scanner.addToken(BANG), nil
		}
	case '=':
		if scanner.match('=') {
			return scanner.addToken(EQUAL_EQUAL), nil
		} else if scanner.match('>') {
			return scanner.addToken(ARROW), nil
		} else {
			return scanner.addToken(EQUAL), nil
		}
	case '<':
		if scanner.match('=') {
			return scanner.addToken(LESS_EQUAL), nil
		} else {
			return scanner.addToken(LESS), nil
		}
	case '>':
		if scanner.match('=') {
			return scanner.addToken(GREATER_EQUAL), nil
		} else {
			return scanner.addToken(GREATER), nil
		}
	case '/':
		if scanner.match('/') {
			return scanner.scanLineComment()
		} else if scanner.match('*') {
			return scanner.scanMultilineComment()
		} else {
			return scanner.addToken(SLASH), nil
		}

		// case '/':
		// 	scanner.token(SLASH)
	}
	return scanner.addToken(ILLEGAL), nil
}

func (scanner *Scanner) scanLineComment() (Token, error) {
	for {
		rune, _, err := scanner.advance()
		if err != nil {
			return scanner.addToken(LINE_COMMENT), err
		}
		if rune == '\n' {
			token := scanner.addToken(LINE_COMMENT)
			scanner.backtrack(rune)
			return token, nil
		}
	}
}

func (scanner *Scanner) scanMultilineComment() (Token, error) {
	for {
		rune, _, err := scanner.advance()
		if err != nil {
			return scanner.addToken(BLOCK_COMMENT), err
		}
		if rune == '*' && scanner.match('/') {
			return scanner.addToken(BLOCK_COMMENT), nil
		}
	}
}

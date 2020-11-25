package main

var scanner Scanner

type Scanner struct {
	Source  string
	Start   int
	Current int
	Line    int
}

func initScanner(source string) {
	scanner.Source = source
	scanner.Start = 0
	scanner.Current = 0
	scanner.Line = 1
}

type Token struct {
	Type   TokenType
	Line   int
	Start  int // N.B. integer offset, not a C-pointer
	Length int
	Source *string
}

func scanToken() Token {
	skipWhitespace()
	scanner.Start = scanner.Current

	if isAtEnd() {
		return makeToken(TOKEN_EOF)
	}
	c := advance()
	if isAlpha(c) {
		return identifier()
	}
	if isDigit(c) {
		return number()
	}
	switch c {
	case '(':
		return makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return makeToken(TOKEN_SEMICOLON)
	case ',':
		return makeToken(TOKEN_COMMA)
	case '.':
		return makeToken(TOKEN_DOT)
	case '-':
		return makeToken(TOKEN_MINUS)
	case '+':
		return makeToken(TOKEN_PLUS)
	case '/':
		return makeToken(TOKEN_SLASH)
	case '*':
		return makeToken(TOKEN_STAR)
	case '!':
		if match('=') {
			return makeToken(TOKEN_BANG_EQUAL)
		} else {
			return makeToken(TOKEN_BANG)
		}
	case '=':
		if match('=') {
			return makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return makeToken(TOKEN_EQUAL)
		}
	case '<':
		if match('=') {
			return makeToken(TOKEN_LESS_EQUAL)
		} else {
			return makeToken(TOKEN_LESS)
		}
	case '>':
		if match('=') {
			return makeToken(TOKEN_GREATER_EQUAL)
		} else {
			return makeToken(TOKEN_GREATER)
		}
	case '"':
		return makeString() // N.B. 'string()' is like a reserved keyword in go.
	}
	return errorToken("unexpected character.")
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func identifier() Token {
	for isAlpha(peek()) || isDigit(peek()) {
		advance()
	}
	return makeToken(identifierType())
}

func identifierType() TokenType {
	switch scanner.Source[scanner.Start] {
	case 'a':
		return checkKeyword(1, 2, "nd", TOKEN_AND)
	case 'c':
		return checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case 'f':
		if scanner.Current-scanner.Start > 1 {
			switch scanner.Source[scanner.Start+1] {
			case 'a':
				return checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return checkKeyword(2, 1, "r", TOKEN_FOR)
			case 'u':
				return checkKeyword(2, 1, "n", TOKEN_FUN)
			}
		}
	case 'i':
		return checkKeyword(1, 1, "f", TOKEN_IF)
	case 'n':
		return checkKeyword(1, 2, "il", TOKEN_NIL)
	case 'o':
		return checkKeyword(1, 1, "r", TOKEN_OR)
	case 'p':
		return checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return checkKeyword(1, 4, "uper", TOKEN_SUPER)
	case 't':
		if scanner.Current-scanner.Start > 1 {
			switch scanner.Source[scanner.Start+1] {
			case 'h':
				return checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	case 'v':
		return checkKeyword(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return checkKeyword(1, 4, "hile", TOKEN_WHILE)
	}
	return TOKEN_IDENTIFIER
}

func checkKeyword(start, length int, rest string, tokenType TokenType) TokenType {
	if scanner.Current-scanner.Start == start+length &&
		string(scanner.Source[scanner.Start+start:scanner.Start+start+length]) == rest {
		return tokenType
	}
	return TOKEN_IDENTIFIER
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func number() Token {
	for isDigit(peek()) {
		advance()
	}
	// check for a fractional part
	if peek() == '.' && isDigit(peekNext()) {
		advance()
		for isDigit(peek()) {
			advance()
		}
	}
	return makeToken(TOKEN_NUMBER)
}

func makeString() Token {
	for peek() != '"' && !isAtEnd() {
		if peek() == '\n' {
			scanner.Line++
		}
		advance()
	}
	if isAtEnd() {
		return errorToken("unterminated string.")
	}
	advance()
	return makeToken(TOKEN_STRING)
}

func skipWhitespace() {
	for {
		c := peek()
		switch c {
		case ' ':
			advance()
		case '\r':
			advance()
		case '\t':
			advance()
		case '\n':
			scanner.Line++
			advance()
		case '/': // skip comments
			if peekNext() == '/' {
				for peek() != '\n' && !isAtEnd() {
					advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func peek() byte {
	if scanner.Current >= len(scanner.Source) { // fake null-terminated strings -.-
		return byte(0)
	}
	return scanner.Source[scanner.Current]
}
func peekNext() byte {
	if isAtEnd() {
		return ' '
	}
	return scanner.Source[scanner.Current+1]
}

func advance() byte {
	scanner.Current++
	return scanner.Source[scanner.Current-1]
}

func match(expected byte) bool {
	if isAtEnd() {
		return false
	}
	if scanner.Source[scanner.Current] != expected {
		return false
	}
	scanner.Current++
	return true
}

func makeToken(tokenType TokenType) Token {
	return Token{
		Type:   tokenType,
		Start:  scanner.Start,
		Length: scanner.Current - scanner.Start,
		Line:   scanner.Line,
		Source: &scanner.Source,
	}
}

func errorToken(message string) Token {
	return Token{
		Type:   TOKEN_ERROR,
		Start:  0,
		Length: len(message),
		Line:   scanner.Line,
		Source: &message,
	}
}

func isAtEnd() bool {
	return scanner.Current >= len(scanner.Source)-1
}

type TokenType byte

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

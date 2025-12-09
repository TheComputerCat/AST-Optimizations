package lexer

type Lexer interface {
	NextToken() Token
}

type lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) Lexer {
	l := &lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.readTwoCharToken(EQ)
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.readTwoCharToken(NEQ)
		} else {
			tok = l.newToken(NOT, l.ch)
		}
	case '(':
		tok = l.newToken(LPAREN, l.ch)
	case ')':
		tok = l.newToken(RPAREN, l.ch)
	case ',':
		tok = l.newToken(COMMA, l.ch)
	case '+':
		tok = l.newToken(PLUS, l.ch)
	case '-':
		tok = l.newToken(MINUS, l.ch)
	case '*':
		tok = l.newToken(ASTARISK, l.ch)
	case '/':
		tok = l.newToken(SLASH, l.ch)
	case '%':
		tok = l.newToken(MODULO, l.ch)
	case '^':
		tok = l.newToken(POWER, l.ch)
	case '~':
		tok = l.newToken(REGEX, l.ch)
	case '<':
		if l.peekChar() == '=' {
			tok = l.readTwoCharToken(LTE)
		} else {
			tok = l.newToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok = l.readTwoCharToken(GTE)
		} else {
			tok = l.newToken(GT, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.readTwoCharToken(AND)
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			tok = l.readTwoCharToken(OR)
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case '"':
		line := l.line
		column := l.column
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isDigit(l.ch) {
			return l.readNumberToken()
		}

		if isLetter(l.ch) {
			line := l.line
			column := l.column
			tok.Literal = l.readIdent()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = column
			return tok
		}

		tok = l.newToken(ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

// Utils to parse input

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// Used for special tokens like <=, >=, !=, etc...
func (l *lexer) readTwoCharToken(tokenType Type) Token {
	ch := l.ch
	line := l.line
	column := l.column
	l.readChar()
	return Token{
		Type:    tokenType,
		Literal: string(ch) + string(l.ch),
		Line:    line,
		Column:  column,
	}
}

func (l *lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// high order function to reuse read + Object (INT, FLOAT, etc...)
func (l *lexer) read(checkFn func(byte) bool) string {
	position := l.position
	for checkFn(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *lexer) readIdent() string {
	return l.read(isLetter)
}

func (l *lexer) readNumber() string {
	return l.read(isDigit)
}

func (l *lexer) readNumberToken() Token {
	line := l.line
	column := l.column
	intPart := l.readNumber()
	if l.ch != '.' {
		return Token{
			Type:    INT,
			Literal: intPart,
			Line:    line,
			Column:  column,
		}
	}

	l.readChar()
	fracPart := l.readNumber()
	return Token{
		Type:    FLOAT,
		Literal: intPart + "." + fracPart,
		Line:    line,
		Column:  column,
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *lexer) newToken(tokenType Type, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}

package lexer

import "minimonkey/token"

type Lexer struct {
	input        string
	position     int  // カーソル位置
	readPosition int  // カーソル位置の次の位置
	ch           byte // カーソル位置の文字
	insertSemi   bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // カーソル位置を初期化する
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	insertSemi := false

	switch ch := l.ch; {
	case isLetter(ch):
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		insertSemi = true
	case isDigit(ch):
		tok.Type = token.INT
		tok.Literal = l.readNumber()
		insertSemi = true
	default:
		switch ch {
		case '=':
			tok = newToken(token.ASSIGN, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
			insertSemi = true
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			if semi := l.insertSemicolon(); semi != nil {
				return *semi
			}
			tok = newToken(token.RBRACE, l.ch)
			insertSemi = true
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case '\n':
			tok = newToken(token.SEMICOLON, ';')
		case 0:
			if semi := l.insertSemicolon(); semi != nil {
				return *semi
			}
			tok = token.Token{Literal: "", Type: token.EOF}
		default:
			tok = newToken(token.ILLEGAL, l.ch)
		}
		l.readChar()
	}

	l.insertSemi = insertSemi

	return tok
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' && !l.insertSemi || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) { // レター文字と数字の組み合わせを許可する
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// // カーソル位置の次の文字を取得する
// func (l *Lexer) peekChar() byte {
// if l.readPosition >= len(l.input) {
// return 0
// }
// return l.input[l.readPosition]
// }

func (l *Lexer) insertSemicolon() *token.Token {
	if l.insertSemi {
		tok := newToken(token.SEMICOLON, ';')
		l.insertSemi = false
		return &tok
	}
	return nil
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

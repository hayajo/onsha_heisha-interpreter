package lexer

import (
	"testing"

	"minimonkey/token"
)

func TestNextToken(t *testing.T) {
	// input := `=+-*/();`
	input := `1 + 2 + 3;
1 + 2 * 3;
(1 + 2) * 3;
let val = 5 + 5;
val + 10;
abc123;
abc123def;
123abc;
123
abc
(1 + 2)`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.PLUS, "+"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.ASTERISK, "*"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.ASTERISK, "*"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "val"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "val"},
		{token.PLUS, "+"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "abc123"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "abc123def"},
		{token.SEMICOLON, ";"},
		{token.INT, "123"},
		{token.IDENT, "abc"},
		{token.SEMICOLON, ";"},
		{token.INT, "123"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "abc"},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] wrong Type. got=%s, expected=%s", i, tok.Type, tt.expectedType)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] wrong Literal. got=%s, expected=%s", i, tok.Literal, tt.expectedLiteral)
		}
	}

}

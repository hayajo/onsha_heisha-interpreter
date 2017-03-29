package lexer

import (
	"testing"

	"minimonkey/token"
)

type tokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

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
(1 + 2)
return`

	tests := []tokenTest{
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

		{token.RETURN, "return"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	testNextToken(t, input, tests)
}

func TestFunctionToken(t *testing.T) {
	input := `let twice = fn(f, x) {
    let once = f(x);
    return f(once);
};

let addTwo = fn(x) { x + 2 }

twice(addTwo, 2)`

	tests := []tokenTest{
		{token.LET, "let"},
		{token.IDENT, "twice"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.COMMA, ","},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.LET, "let"},
		{token.IDENT, "once"},
		{token.ASSIGN, "="},
		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.RETURN, "return"},
		{token.IDENT, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "once"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "addTwo"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "twice"},
		{token.LPAREN, "("},
		{token.IDENT, "addTwo"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	testNextToken(t, input, tests)
}

func testNextToken(t *testing.T, input string, tests []tokenTest) {
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

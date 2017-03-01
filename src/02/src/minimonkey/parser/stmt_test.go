package parser

import (
	"minimonkey/ast"
	"minimonkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = x;", "y", "x"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.Parse()

		checkParseErrors(t, p)
		if program == nil {
			t.Fatalf("Parse() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 3)
		}

		stmt := program.Statements[0]
		if stmt.TokenLiteral() != "let" {
			t.Errorf("stmt.TokenLiteral() got %s, expected %s", stmt.TokenLiteral(), "let")
			return
		}

		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			t.Errorf("stmt is %s, expect %T", letStmt, "*ast.LetStatement")
			return
		}

		switch v := tt.expectedValue.(type) {
		case int:
			if !testIntegerLiteral(t, letStmt.Value, int64(v)) {
				return
			}
		case int64:
			if !testIntegerLiteral(t, letStmt.Value, v) {
				return
			}
		case string:
			if !testIdentifier(t, letStmt.Value, v) {
				return
			}
		}
	}

}

func TestInvalidStatement(t *testing.T) {
	tt := []struct {
		input string
	}{
		{input: "let x;"},
		{input: "let x = 3 3 - 3;"},
		{input: "5 let x = 3;"},
		{input: "x y z;"},
		{input: "x 5 z;"},
		{input: "1 2 3;"},
		{input: "123abc;"},
	}

	for _, tt := range tt {
		l := lexer.New(tt.input)
		p := New(l)

		p.Parse()

		if len(p.Errors()) == 0 {
			t.Fatalf("p.Parse() is expected to be error")
		}
	}
}

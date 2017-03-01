package parser

import (
	"minimonkey/ast"
	"minimonkey/lexer"
	"strconv"
	"testing"
)

func TestIdentifier(t *testing.T) {
	input := `foobar`

	l := lexer.New(input)
	p := New(l)

	program := p.Parse()

	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 1)
	}

	if program.Statements[0].TokenLiteral() != "foobar" {
		t.Errorf("stmt.TokenLiteral() got %s, expected %s", program.Statements[0].TokenLiteral(), "foobar")
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is %T, expect %s", stmt, "*ast.ExpressionStatement")
		return
	}

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("ast.Expression is %T, expect %s", exp, "*ast.Identifier")
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is %d, expect %d", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is %s, expect %s", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func TestIntegerLiteral(t *testing.T) {
	input := `5`

	l := lexer.New(input)
	p := New(l)

	program := p.Parse()

	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 1)
	}

	if program.Statements[0].TokenLiteral() != "5" {
		t.Errorf("stmt.TokenLiteral() got %s, expected %s", program.Statements[0].TokenLiteral(), "5")
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is %T, expect %s", stmt, "*ast.ExpressionStatement")
		return
	}

	il, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("stmt.Expression is %T, expect %s", il, "*ast.IntegerLiteral")
		return
	}
	if il.Value != 5 {
		t.Errorf("il.Value is %d, expect %d", il.Value, 5)
		return
	}
	if il.TokenLiteral() != "5" {
		t.Errorf("stmt.TokenLiteral() is %s, expect %s", il.TokenLiteral(), "5")
		return
	}
}

func TestParsingPrefixOperator(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"-15", "-", 15},
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
			t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is %T, expect %s", stmt, "*ast.ExpressionStatement")
			return
		}

		pexp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt.Expression is %T, expect %s", pexp, "*ast.PrefixExpression")
			return
		}
		if pexp.Operator != tt.operator {
			t.Errorf("pexp.Operator is %s, expect %s", pexp.Operator, tt.operator)
			return
		}
		if !testIntegerLiteral(t, pexp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("ast.Expression is %T, expect %s", exp, "*ast.IntegerLiteral")
		return false
	}

	if il.Value != value {
		t.Errorf("il.Value is %d, expect %d", il.Value, value)
		return false
	}

	if il.TokenLiteral() != strconv.FormatInt(value, 10) {
		t.Errorf("il.TokenLiteral() is %s, expect %s", il.TokenLiteral(), strconv.FormatInt(value, 10))
		return false
	}

	return true
}

func TestParsingInfixOperator(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
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
			t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is %T, expect %s", stmt, "*ast.ExpressionStatement")
			return
		}

		iexp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("stmt.Expression is %T, expect %s", iexp, "*ast.InfixExpression")
			return
		}
		if !testIntegerLiteral(t, iexp.Left, tt.leftValue) {
			return
		}
		if iexp.Operator != tt.operator {
			t.Errorf("iexp.Operator is %s, expect %s", iexp.Operator, tt.operator)
			return
		}
		if !testIntegerLiteral(t, iexp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input   string
		expectd string
	}{
		{"-a * b", "((-a) * b);"},
		{"a + b + c", "((a + b) + c);"},
		{"a + b - c", "((a + b) - c);"},
		{"a * b * c", "((a * b) * c);"},
		{"a * b / c", "((a * b) / c);"},
		{"a + b / c", "(a + (b / c));"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f);"},
		{"3 + 4; -5 * 5", "(3 + 4);((-5) * 5);"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4);"},
		{"(5 + 5) * 2", "((5 + 5) * 2);"},
		{"2 / (5 + 5)", "(2 / (5 + 5));"},
		{"-(5 + 5)", "(-(5 + 5));"},
		{"6 / 2 * (1 + 2)", "((6 / 2) * (1 + 2));"},
		{"6 / (2 * (1 + 2))", "(6 / (2 * (1 + 2)));"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.Parse()

		checkParseErrors(t, p)
		if program == nil {
			t.Fatalf("Parse() returned nil")
		}

		got := program.String()
		if got != tt.expectd {
			t.Fatalf("program.String() got %q, expected %q", got, tt.expectd)
		}
	}
}

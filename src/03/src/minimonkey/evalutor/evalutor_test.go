package evalutor

import (
	"minimonkey/lexer"
	"minimonkey/object"
	"minimonkey/parser"
	"testing"
)

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestEvalLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)

		errObj, ok := evaluted.(*object.Error)
		if !ok {
			t.Errorf("object is not Error got %T (%+v)", evaluted, evaluted)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("errObj.Message got %q, expected %q", errObj.Message, tt.expected)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	v, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer got %T (%+v)", obj, obj)
		return false
	}
	if v.Value != expected {
		t.Errorf("object has wrong value. got %d, expected %d", v.Value, expected)
		return false
	}
	return true
}

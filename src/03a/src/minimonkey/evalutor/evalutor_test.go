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
		{"let v = fn(){}(); v + 1;", "unknown operator NULL + INTEGER"},
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

func TestEvalEmptyStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{";", nil},
		{"; 10;", 10},
		{"10; ;", nil},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)

		if tt.expected == nil {
			if evaluted != NULL {
				t.Error("evaluted is not a NULL")
			}
		} else {
			testIntegerObject(t, evaluted, int64(tt.expected.(int)))
		}
	}
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"return;", nil},
		{"return; return 10", nil},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)

		if tt.expected == nil {
			if evaluted != NULL {
				t.Error("evaluted is not a NULL")
			}
		} else {
			testIntegerObject(t, evaluted, int64(tt.expected.(int)))
		}
	}
}

func TestEvalFunctionLiteral(t *testing.T) {
	tests := []struct {
		input  string
		params []string
		body   string
	}{
		{"fn() { }", []string{}, "{}"},
		{"fn() { return }", []string{}, "{ return; }"},
		{"fn(x) { x + 2 }", []string{"x"}, "{ (x + 2); }"},
		{"fn(x, y) { return x + y + 2; }", []string{"x", "y"}, "{ return ((x + y) + 2); }"},
		{"fn(x) { x + 1; x + 2 }", []string{"x"}, "{ (x + 1); (x + 2); }"},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		fn, ok := evaluted.(*object.Function)
		if !ok {
			t.Errorf("evaluted is %T, expected %s", evaluted, "Function")
			return
		}

		if len(fn.Parameters) != len(tt.params) {
			t.Errorf("fn.Parameters has wrong parametes. expect %d parameters", len(tt.params))
			return
		}

		for i, p := range fn.Parameters {
			if p.String() != tt.params[i] {
				t.Errorf("p.String() got %q, expect %q", p.String(), tt.params[i])
				return
			}
		}

		if fn.Body.String() != tt.body {
			t.Errorf("fn.Body.String() got %q, expect %q", fn.Body.String(), tt.body)
			return
		}
	}
}

func TestEvalCallExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let identity = fn(x) { x }; identity(5);", 5},
		{"let identity = fn(x) { return x }; identity(5);", 5},
		{"let double = fn(x) { x * 2 }; double(5);", 10},
		{"let add = fn(x, y) { x + y }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y }; add(5 + 5, add(5, 5));", 20},
		{"let x = 5; let y = fn(x) { x }(10); x;", 5},
		{"fn(x) { x }(5)", 5},
		{"fn() {}()", nil},
		{"fn() { return }()", nil},
		{"let callTwoTimes = fn(x, func) { func(func(x)) }; callTwoTimes(3, fn(x) { x + 1 });", 5}, // high-oder function
		{"let newAdder = fn(x) { fn(n) { x + n } }; let addTwo = newAdder(2); addTwo(2);", 4},      // closure
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)

		if tt.expected == nil {
			if evaluted != NULL {
				t.Error("evaluted is not a NULL")
			}
		} else {
			testIntegerObject(t, evaluted, int64(tt.expected.(int)))
		}
	}
}

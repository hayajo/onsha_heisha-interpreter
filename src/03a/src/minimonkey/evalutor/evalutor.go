package evalutor

import (
	"fmt"

	"minimonkey/ast"
	"minimonkey/object"
)

var (
	NULL = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)

	// // TODO
	// case *ast.EmptyStatement:
	// case *ast.ReturnStatement:
	// case *ast.BlockStatement:

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: %s", node.Value)
		}
		return val

		// // TODO
		// case *ast.FunctionLiteral
		// case *ast.CallExpression
	}

	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var res object.Object

	for _, s := range program.Statements {
		res = Eval(s, env)
	}

	return res
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(val object.Object) object.Object {
	if val.Type() != object.INTEGER_OBJ {
		return newError("unknown operator -%s", val.Type())
	}

	v := val.(*object.Integer).Value

	return &object.Integer{Value: -v}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ || right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
	}

	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	default:
		return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

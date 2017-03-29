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
		// return evalProgram(node, env)
		return evalStatements(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)

	case *ast.EmptyStatement:
		return NULL

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if val == nil {
			val = NULL
		}
		return &object.ReturnValue{Value: val}

	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)

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

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		val := applyFunction(function, args)
		if val == nil {
			return NULL
		}

		return val
	}

	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// func evalProgram(program *ast.Program, env *object.Environment) object.Object {
func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var res object.Object

	for _, s := range statements {
		res = Eval(s, env)

		if v, ok := res.(*object.ReturnValue); ok {
			return v.Value
		}
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

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(exps))

	for _, exp := range exps {
		evaluted := Eval(exp, env)
		if isError(evaluted) {
			return []object.Object{evaluted}
		}
		result = append(result, evaluted)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	env := extendFunctionEnv(function, args) // TODO: don't yet support closure
	evaluted := Eval(function.Body, env)

	if returnValue, ok := evaluted.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluted
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

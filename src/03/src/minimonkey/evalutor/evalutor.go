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
	// // TODO
	// // Statements
	// case *ast.Program:
	// case *ast.ExpressionStatement:
	// case *ast.LetStatement:
	// // Expressions
	// case *ast.IntegerLiteral:
	// case *ast.PrefixExpression:
	// case *ast.InfixExpression:
	// case *ast.Identifier:
	}
	return newError("Unknown node type %T", node)
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

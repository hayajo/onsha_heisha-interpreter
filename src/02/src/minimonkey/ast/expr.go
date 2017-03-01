package ast

import (
	"bytes"
	"minimonkey/token"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pexp *PrefixExpression) expressionNode()      {}
func (pexp *PrefixExpression) TokenLiteral() string { return pexp.Token.Literal }
func (pexp *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pexp.Operator)
	out.WriteString(pexp.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (iexp *InfixExpression) expressionNode()      {}
func (iexp *InfixExpression) TokenLiteral() string { return iexp.Token.Literal }
func (iexp *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(iexp.Left.String())
	out.WriteString(" " + iexp.Operator + " ")
	out.WriteString(iexp.Right.String())
	out.WriteString(")")
	return out.String()
}

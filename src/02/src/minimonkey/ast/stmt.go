package ast

import (
	"bytes"
	"minimonkey/token"
)

type Statement interface {
	Node
	statementNode()
}

type EmptyStatement struct {
	Token token.Token
}

func (es *EmptyStatement) statementNode()       {}
func (es *EmptyStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EmptyStatement) String() string       { return "" }

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
		out.WriteString(";")
	}

	return out.String()
}

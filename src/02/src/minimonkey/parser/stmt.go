package parser

import (
	"minimonkey/ast"
	"minimonkey/token"
)

func (p *Parser) parseStmt() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.SEMICOLON:
		return &ast.EmptyStatement{Token: p.curToken}, nil
	default:
		return p.parseExpressionStatement()
	}
}

// let <identifier> = <expression>;
func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil, p.peekError(token.IDENT)
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil, p.peekError(token.ASSIGN)
	}

	p.nextToken()

	value, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Value = value

	if !p.expectPeek(token.SEMICOLON) {
		return nil, p.peekError(token.SEMICOLON)
	}

	return stmt, nil
}

// <expression>;
func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = exp

	if !p.expectPeek(token.SEMICOLON) {
		return nil, p.peekError(token.SEMICOLON)
	}

	return stmt, nil
}

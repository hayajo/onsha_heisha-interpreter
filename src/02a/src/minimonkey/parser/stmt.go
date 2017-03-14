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
	case token.RETURN:
		return p.parseReturnStatement()
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

// return [expression];
func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	if !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		v, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.ReturnValue = v
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil, p.peekError(token.SEMICOLON)
	}

	return stmt, nil
}

// { [statement...] }
func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken} // curToken == LBRACE
	block.Statements = []ast.Statement{}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
	}

	if !p.expectPeek(token.RBRACE) {
		return nil, p.peekError(token.RBRACE)
	}

	return block, nil
}

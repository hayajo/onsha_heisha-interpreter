package parser

import (
	"minimonkey/ast"
	"minimonkey/lexer"
	"minimonkey/token"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		expect *ast.Program
	}{
		{
			input: "-1 + (2 * 3 - 4) * 5",
			expect: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{
							Literal: "+",
							Type:    token.PLUS,
						},
						Expression: &ast.InfixExpression{
							Operator: "+",
							Token: token.Token{
								Literal: "+",
								Type:    token.PLUS,
							},
							Left: &ast.PrefixExpression{
								Operator: "-",
								Token: token.Token{
									Literal: "-",
									Type:    token.MINUS,
								},
								Right: &ast.IntegerLiteral{
									Token: token.Token{
										Literal: "1",
										Type:    token.INT,
									},
									Value: 1,
								},
							},
							Right: &ast.InfixExpression{
								Operator: "*",
								Token: token.Token{
									Literal: "*",
									Type:    token.ASTERISK,
								},
								Left: &ast.InfixExpression{
									Operator: "-",
									Token: token.Token{
										Literal: "-",
										Type:    token.MINUS,
									},
									Left: &ast.InfixExpression{
										Operator: "*",
										Token: token.Token{
											Literal: "*",
											Type:    token.ASTERISK,
										},
										Left: &ast.IntegerLiteral{
											Token: token.Token{
												Literal: "2",
												Type:    token.INT,
											},
											Value: 2,
										},
										Right: &ast.IntegerLiteral{
											Token: token.Token{
												Literal: "3",
												Type:    token.INT,
											},
											Value: 3,
										},
									},
									Right: &ast.IntegerLiteral{
										Token: token.Token{
											Literal: "4",
											Type:    token.INT,
										},
										Value: 4,
									},
								},
								Right: &ast.IntegerLiteral{
									Token: token.Token{
										Literal: "5",
										Type:    token.INT,
									},
									Value: 5,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		p := New(l)
		program := p.Parse()

		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements contain %d statements, expected %d", len(program.Statements), 1)
		}

		if program.String() != tt.expect.String() {
			t.Errorf("program.String() got %q, expect %q", program.String(), tt.expect.String())
		}
	}

}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

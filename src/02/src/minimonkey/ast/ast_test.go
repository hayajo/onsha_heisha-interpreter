package ast

import "minimonkey/token"
import "testing"

func TestString(t *testing.T) {
	tests := []struct {
		program *Program
		expect  string
	}{
		{
			&Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{
							Type:    token.LET,
							Literal: "let",
						},
						Name: &Identifier{
							Token: token.Token{
								Type:    token.LET,
								Literal: "muVar",
							},
							Value: "myVar",
						},
						Value: &Identifier{
							Token: token.Token{
								Type:    token.LET,
								Literal: "anotherVar",
							},
							Value: "anotherVar",
						},
					},
				},
			},
			"let myVar = anotherVar;",
		},
		{
			&Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{
							Type:    token.LET,
							Literal: "let",
						},
						Name: &Identifier{
							Token: token.Token{
								Type:    token.LET,
								Literal: "muVar",
							},
							Value: "myVar",
						},
						Value: &Identifier{
							Token: token.Token{
								Type:    token.INT,
								Literal: "10",
							},
							Value: "10",
						},
					},
				},
			},
			"let myVar = 10;",
		},
	}

	for _, tt := range tests {
		if tt.program.String() != tt.expect {
			t.Errorf("program.String() got %q, expect %q", tt.program.String(), tt.expect)
		}
	}
}

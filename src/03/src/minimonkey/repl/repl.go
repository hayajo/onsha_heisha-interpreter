package repl

import (
	"bufio"
	"fmt"
	"io"

	"minimonkey/evalutor"
	"minimonkey/lexer"
	"minimonkey/object"
	"minimonkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.Parse()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		evaluted := evalutor.Eval(program, env)

		io.WriteString(out, evaluted.Inspect())
		io.WriteString(out, "\n")
	}
}

func printErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "ERROR: "+err.Error()+"\n")
	}
}

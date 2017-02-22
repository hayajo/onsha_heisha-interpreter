package main

import (
	"fmt"
	"minimonkey/repl"
	"os"
)

func main() {
	fmt.Println("This is the MiniMonkey programming language!\n")
	repl.Start(os.Stdin, os.Stdout)
}

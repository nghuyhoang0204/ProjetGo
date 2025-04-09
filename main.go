package main

import (
	"ProjetGo/codegen"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"fmt"
)

func main() {
	// Input TypeScript code (can be read from a file or hardcoded)
	input := `let x: number = 10;`

	// Step 1: Tokenize the input
	tokens := lexer.Tokenize(input)

	// Step 2: Parse the tokens and build the AST
	ast := parser.Parse(tokens)

	// Step 3: Generate JavaScript code from the AST
	jsCode := codegen.GenerateCode(ast)

	// Output the generated JavaScript code
	fmt.Println("Generated JavaScript Code:")
	fmt.Println(jsCode)
}

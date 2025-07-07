package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func testProgressiveFeatures() {
	// Test 1: Variables (déjà testé)
	fmt.Println("=== Test 1: Variables ===")
	test1 := `
const name: string = "Alice";
let score: number = 100;
let isWinner: boolean = true;
`
	runTest(test1)

	// Test 2: Expressions simples
	fmt.Println("\n=== Test 2: Expression simple ===")
	test2 := `
let result = 5 + 3;
`
	runTest(test2)

	// Test 3: Return statement
	fmt.Println("\n=== Test 3: Return statement ===")
	test3 := `
return 42;
`
	runTest(test3)
}

func runTest(code string) {
	fmt.Println("Code:", code)
	
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("Parsé: %d éléments\n", len(program))
	
	for i, stmt := range program {
		if stmt != nil {
			fmt.Printf("  %d. %T\n", i+1, stmt)
		}
	}

	jsOutput := generator.Generate(program, generator.JavaScript)
	fmt.Println("JavaScript:", jsOutput)
}

func main() {
	testProgressiveFeatures()
}

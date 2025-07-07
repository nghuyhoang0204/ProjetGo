package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func testSimplifiedCode() {
	// Test avec du code simplifié d'abord
	code := `// Déclaration de variables
const nom: string = "Alice";
let age: number = 25;
var actif: boolean = true;

// Types de base
let score: number = 42;
let message: string = "Hello World";
let valeurs: number[] = [10, 20, 30];`

	fmt.Println("=== Code Test Simplifié ===")
	fmt.Println(code)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("\nÉléments parsés: %d\n", len(program))
	
	for i, stmt := range program {
		if stmt != nil {
			fmt.Printf("%d. %T\n", i+1, stmt)
		}
	}

	// Test JavaScript
	jsOutput := generator.Generate(program, generator.JavaScript)
	fmt.Println("\n=== JavaScript ===")
	fmt.Println(jsOutput)

	// Test Java
	javaOutput := generator.Generate(program, generator.Java)
	fmt.Println("=== Java ===")
	fmt.Println(javaOutput)
}

func main() {
	testSimplifiedCode()
}

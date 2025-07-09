package main

import (
	"ProjetGo/generator"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"fmt"
	"testing"
)

// Test simple de la transpilation TypeScript vers JavaScript
func TestTypescriptToJavascript(t *testing.T) {
	// Code TypeScript à tester
	input := `
function addition(a: number, b: number): number {
  return a + b;
}
const resultat = addition(5, 3);
console.log("Résultat :", resultat); // Résultat : 8
`

	// Transpilation
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	
	// Génération JavaScript
	output := generator.Generate(program.Statements)

	// Afficher le résultat
	fmt.Println("TypeScript original:")
	fmt.Println("-------------------")
	fmt.Println(input)
	fmt.Println("\nJavaScript généré:")
	fmt.Println("------------------")
	fmt.Println(output)
}

package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func debugParsing() {
	// Code TypeScript simple pour test progressif
	simpleCode := `
const nom: string = "Alice";
let age: number = 25;
const numbers: number[] = [10, 20, 30];
const person = {
  name: "Bob",
  age: 30
};
console.log("Hello World");
console.log(nom);
`

	fmt.Println("=== Test avec code TypeScript simple ===")
	fmt.Println(simpleCode)

	l := lexer.New(simpleCode)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("Éléments parsés: %d\n", len(program))

	// Afficher chaque élément parsé
	for i, stmt := range program {
		if stmt != nil {
			fmt.Printf("%d. Type: %T, Token: %s\n", i+1, stmt, stmt.TokenLiteral())
		}
	}

	// Test de génération pour tous les langages
	fmt.Println("\n=== JavaScript Output ===")
	jsOutput := generator.Generate(program, generator.JavaScript)
	fmt.Println(jsOutput)
	
	fmt.Println("\n=== Java Output ===")
	javaOutput := generator.Generate(program, generator.Java)
	fmt.Println(javaOutput)
	
	fmt.Println("\n=== Python Output ===")
	pythonOutput := generator.Generate(program, generator.Python)
	fmt.Println(pythonOutput)
}

func main() {
	debugParsing()
}

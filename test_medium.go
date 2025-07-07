package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func testSimpleToMediumCode() {
	// Code simple à moyen complexe
	code := `
// Variables
const name: string = "John Doe";
let age: number = 25;
let isActive: boolean = true;

// Fonction simple
function greet(name: string): string {
    return "Hello " + name;
}

// Condition
if (age >= 18) {
    console.log("Adult");
} else {
    console.log("Minor");
}

// Boucle for
for (let i = 0; i < 5; i++) {
    console.log(i);
}

// Boucle while
let counter = 0;
while (counter < 3) {
    counter++;
}
`

	fmt.Println("=== Code TypeScript/JavaScript ===")
	fmt.Println(code)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("\nÉléments parsés: %d\n", len(program))

	// Test tous les langages
	languages := map[string]generator.TargetLanguage{
		"JavaScript": generator.JavaScript,
		"Java":       generator.Java,
		"Python":     generator.Python,
		"C#":         generator.CSharp,
		"Go":         generator.Go,
	}

	for name, lang := range languages {
		output := generator.Generate(program, lang)
		fmt.Printf("\n=== %s Output ===\n", name)
		fmt.Println(output)
	}
}

func main() {
	testSimpleToMediumCode()
}

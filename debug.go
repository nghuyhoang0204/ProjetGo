package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func debugParsing() {
	// Votre exemple TypeScript → JavaScript
	tsCode := `
const nom: string = "Lucie";
let age: number = 17;
var majeur: boolean = false;

function saluer(n: string): void {
  console.log("Bonjour " + n);
}

if (age >= 18) {
  majeur = true;
} else {
  majeur = false;
}

let notes: number[] = [12, 15, 9];
let eleve = { nom: nom, age: age };

for (let i = 0; i < notes.length; i++) {
  console.log("Note :", notes[i]);
}

let compteur: number = 3;
while (compteur > 0) {
  console.log("Compte :", compteur);
  compteur--;
}

saluer(eleve.nom);
console.log("Est majeur :", majeur);
`

	fmt.Println("=== TypeScript Input ===")
	fmt.Println(tsCode)
	
	// Parser et générer
	l := lexer.New(tsCode)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("Éléments parsés: %d\n", len(program))

	// Afficher chaque élément parsé avec plus de détails
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
	
	fmt.Println("\n=== C# Output ===")
	csharpOutput := generator.Generate(program, generator.CSharp)
	fmt.Println(csharpOutput)
	
	fmt.Println("\n=== Go Output ===")
	goOutput := generator.Generate(program, generator.Go)
	fmt.Println(goOutput)
}

func main() {
	debugParsing()
}

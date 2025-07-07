package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
)

func testUserCode() {
	// Votre code exact
	code := `// Déclaration de variables
const nom: string = "Alice";         // constante
let age: number = 25;                // variable modifiable
var actif: boolean = true;           // ancienne syntaxe

// Types de base
let score: number = 42;
let message: string = "Hello World";
let valeurs: number[] = [10, 20, 30];
let personne: { nom: string; age: number } = {
  nom: "Bob",
  age: 30
};

// Fonction simple
function saluer(prenom: string): string {
  return "Salut !";
}

// Condition simple
function verifierAge(age: number): string {
  if (age >= 18) {
    return "Majeur";
  } else {
    return "Mineur";
  }
}

// Utilisation
console.log(message);
console.log(saluer("Marie"));`

	fmt.Println("=== Code TypeScript Complet ===")
	fmt.Println(code)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Printf("\nÉléments parsés: %d\n", len(program))
	
	for i, stmt := range program {
		if stmt != nil {
			fmt.Printf("%d. %T - %s\n", i+1, stmt, stmt.TokenLiteral())
		}
	}

	// Test des 5 langages
	languages := map[string]generator.TargetLanguage{
		"JavaScript": generator.JavaScript,
		"Java":       generator.Java,
		"Python":     generator.Python,
		"C#":         generator.CSharp,
		"Go":         generator.Go,
	}

	for name, lang := range languages {
		output := generator.Generate(program, lang)
		fmt.Printf("\n=== %s ===\n", name)
		fmt.Println(output)
	}
}

func main() {
	testUserCode()
}

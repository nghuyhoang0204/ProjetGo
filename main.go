package main

import (
	"ProjetGo/lexer"
    "ProjetGo/parser"
    "ProjetGo/generator"
    "fmt"
    "os"
)

func runConsoleVersion() {
	input := `
	const message: string = "Hello World";
	let count: number = 42;
	const pi: number = 3.14;
	`

	fmt.Println("=== Code Source (TypeScript-like) ===")
	fmt.Println(input)

	l := lexer.New(input)         
	p := parser.New(l)            
	program := p.ParseProgram() 

	// Génération en JavaScript
	jsOutput := generator.Generate(program, generator.JavaScript)
	fmt.Println("\n=== JavaScript Output ===")
	fmt.Println(jsOutput)

	// Génération en Java
	javaOutput := generator.Generate(program, generator.Java)
	fmt.Println("=== Java Output ===")
	fmt.Println(javaOutput)

	// Génération en Python
	pythonOutput := generator.Generate(program, generator.Python)
	fmt.Println("=== Python Output ===")
	fmt.Println(pythonOutput)

	// Génération en C#
	csharpOutput := generator.Generate(program, generator.CSharp)
	fmt.Println("=== C# Output ===")
	fmt.Println(csharpOutput)

	// Génération en Go
	goOutput := generator.Generate(program, generator.Go)
	fmt.Println("=== Go Output ===")
	fmt.Println(goOutput)
}

func main() {
	fmt.Println("🚀 Transpilateur Multi-Langages")
	fmt.Println("================================")
	
	// Vérifier les arguments de ligne de commande
	if len(os.Args) > 1 && os.Args[1] == "console" {
		fmt.Println("Mode console activé...")
		runConsoleVersion()
		return
	}

	// Par défaut, lancer l'interface web
	fmt.Println("🌐 Lancement de l'interface web...")
	fmt.Println("💡 Pour utiliser la version console: go run . console")
	fmt.Println("")
	
	StartWebServer()
}

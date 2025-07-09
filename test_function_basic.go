package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
	"os"
)

func testBasicFunction() {
	fmt.Println("=== Test Fonction Basique ===")
	
	// Test étape par étape
	tests := []struct {
		name string
		code string
	}{
		{
			name: "Variable simple",
			code: `let age: number = 17;`,
		},
		{
			name: "Fonction vide",
			code: `function saluer(): void { }`,
		},
		{
			name: "Fonction avec paramètre",
			code: `function saluer(n: string): void { }`,
		},
		{
			name: "Fonction avec console.log",
			code: `function saluer(n: string): void {
  console.log("Bonjour");
}`,
		},
		{
			name: "If statement simple",
			code: `if (age >= 18) {
  console.log("Majeur");
}`,
		},
	}

	for _, test := range tests {
		fmt.Printf("\n--- %s ---\n", test.name)
		fmt.Printf("Code: %s\n", test.code)
		
		// Essayer de parser
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("❌ CRASH: %v\n", r)
				}
			}()
			
			l := lexer.New(test.code)
			p := parser.New(l)
			program := p.ParseProgram()
			
			fmt.Printf("✅ Parsing OK - %d éléments\n", len(program))
			
			// Test génération JavaScript
			if len(program) > 0 {
				output := generator.Generate(program, generator.JavaScript)
				fmt.Printf("JavaScript: %s\n", output)
			}
		}()
	}
}

func testComplexFunctionSafe(code string, targetLang string) {
	fmt.Printf("=== Test Fonction Complexe pour %s ===\n", targetLang)
	fmt.Printf("Code:\n%s\n\n", code)
	
	filename := fmt.Sprintf("test_parsing_%s.txt", targetLang)
	
	// Supprimer le fichier précédent
	os.Remove(filename)
	
	defer func() {
		if r := recover(); r != nil {
			errorContent := fmt.Sprintf("PARSING ERROR\n=============\n\nCode:\n%s\n\nTarget: %s\n\nError: %v", code, targetLang, r)
			os.WriteFile(filename, []byte(errorContent), 0644)
			fmt.Printf("❌ CRASH sauvegardé dans %s\n", filename)
		}
	}()
	
	// Parser
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	
	fmt.Printf("✅ Parsing réussi - %d éléments\n", len(program))
	
	// Afficher les éléments parsés
	for i, stmt := range program {
		if stmt != nil {
			fmt.Printf("%d. Type: %T\n", i+1, stmt)
		}
	}
	
	// Générer selon le langage
	var lang generator.TargetLanguage
	switch targetLang {
	case "javascript":
		lang = generator.JavaScript
	case "java":
		lang = generator.Java
	case "python":
		lang = generator.Python
	default:
		lang = generator.JavaScript
	}
	
	output := generator.Generate(program, lang)
	
	// Créer le contenu du fichier
	content := fmt.Sprintf(`PARSING TEST RESULTS
====================

SOURCE CODE:
%s

TARGET LANGUAGE: %s
PARSED ELEMENTS: %d

ELEMENTS DETAILS:
`, code, targetLang, len(program))
	
	for i, stmt := range program {
		if stmt != nil {
			content += fmt.Sprintf("%d. Type: %T, Token: %s\n", i+1, stmt, stmt.TokenLiteral())
		}
	}
	
	content += fmt.Sprintf(`

GENERATED %s CODE:
=====================
%s
`, targetLang, output)
	
	// Sauvegarder
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("❌ Erreur sauvegarde: %v\n", err)
	} else {
		fmt.Printf("✅ Résultat sauvegardé dans %s\n", filename)
	}
	
	// Afficher un aperçu du résultat
	fmt.Printf("\nAperçu du code généré:\n%s\n", output)
}

func main() {
	// Test basique d'abord
	testBasicFunction()
	
	fmt.Println("\n" + "="*60 + "\n")
	
	// Test fonction complexe
	complexCode := `function saluer(n: string): void {
  console.log("Bonjour " + n);
}`
	
	// Tester pour JavaScript
	testComplexFunctionSafe(complexCode, "javascript")
	
	fmt.Println("\n" + "="*30 + "\n")
	
	// Tester pour Java
	testComplexFunctionSafe(complexCode, "java")
	
	fmt.Println("\n" + "="*30 + "\n")
	
	// Code plus complexe
	moreComplexCode := `let age: number = 17;

function saluer(n: string): void {
  console.log("Bonjour " + n);
}

if (age >= 18) {
  console.log("Majeur");
} else {
  console.log("Mineur");
}`
	
	// Tester le code plus complexe
	testComplexFunctionSafe(moreComplexCode, "javascript")
}

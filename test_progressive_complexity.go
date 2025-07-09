package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
	"os"
)

func testProgressiveComplexity() {
	fmt.Println("=== TEST PROGRESSIF DE COMPLEXITÉ ===")
	
	// Test 1: Variables + console.log (fonctionne)
	test1 := `const nom: string = "Alice";
console.log("Bonjour " + nom);`
	
	// Test 2: Variables + console.log + function simple
	test2 := `const nom: string = "Alice";
function saluer(): void {
  console.log("Bonjour");
}`
	
	// Test 3: Variables + function avec paramètres
	test3 := `const nom: string = "Alice";
function saluer(n: string): void {
  console.log("Bonjour " + n);
}`
	
	// Test 4: Variables + function + if/else
	test4 := `const nom: string = "Alice";
let age: number = 25;
function saluer(n: string): void {
  console.log("Bonjour " + n);
}
if (age >= 18) {
  console.log("Majeur");
}`
	
	// Test 5: Ajout des arrays
	test5 := `const nom: string = "Alice";
let notes: number[] = [12, 15, 9];`
	
	// Test 6: Ajout des objets
	test6 := `const nom: string = "Alice";
let eleve = { nom: nom, age: 25 };`
	
	// Test 7: Ajout for loop
	test7 := `const nom: string = "Alice";
let notes: number[] = [12, 15, 9];
for (let i = 0; i < notes.length; i++) {
  console.log("Note " + i);
}`
	
	tests := []struct {
		name string
		code string
	}{
		{"Variables + console.log", test1},
		{"+ Function simple", test2},
		{"+ Function avec paramètres", test3},
		{"+ Function + if/else", test4},
		{"+ Arrays", test5},
		{"+ Objets", test6},
		{"+ For loop", test7},
	}
	
	for i, test := range tests {
		fmt.Printf("\n🧪 Test %d: %s\n", i+1, test.name)
		fmt.Println("Code:", test.code)
		
		success := testParsing(test.code, fmt.Sprintf("test_%d_result.txt", i+1))
		if success {
			fmt.Printf("✅ Test %d réussi\n", i+1)
		} else {
			fmt.Printf("❌ Test %d échoué - ARRÊT ICI\n", i+1)
			fmt.Printf("🔍 Le problème vient de: %s\n", test.name)
			break
		}
	}
}

func testParsing(code string, filename string) bool {
	defer func() {
		if r := recover(); r != nil {
			content := fmt.Sprintf("CRASH DETECTED\n=============\n\nCode: %s\n\nError: %v", code, r)
			os.WriteFile(filename, []byte(content), 0644)
		}
	}()
	
	// Tentative de parsing
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	
	if len(program) == 0 {
		content := fmt.Sprintf("NO PARSING\n==========\n\nCode: %s\n\nParsed Elements: 0", code)
		os.WriteFile(filename, []byte(content), 0644)
		return false
	}
	
	// Tentative de génération JavaScript
	output := generator.Generate(program, generator.JavaScript)
	
	content := fmt.Sprintf("SUCCESS\n=======\n\nCode: %s\n\nParsed Elements: %d\n\nJavaScript Output:\n%s", 
		code, len(program), output)
	os.WriteFile(filename, []byte(content), 0644)
	
	return true
}

func main() {
	testProgressiveComplexity()
}

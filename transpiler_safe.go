package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

// Structure pour sauvegarder le parsing
type SavedParsing struct {
	SourceCode    string                 `json:"source_code"`
	ElementsCount int                    `json:"elements_count"`
	ParsedOK      bool                   `json:"parsed_ok"`
	ErrorMessage  string                 `json:"error_message"`
	ParsedTypes   []string               `json:"parsed_types"`
}

// Sauvegarder le parsing UNE SEULE FOIS (sans affichage terminal)
func saveParsingOnce(code string, filename string) {
	fmt.Printf("🔄 Parsing du code en cours (sauvegarde dans %s)...\n", filename)
	
	saved := SavedParsing{
		SourceCode: code,
		ParsedOK:   false,
	}

	defer func() {
		if r := recover(); r != nil {
			saved.ErrorMessage = fmt.Sprintf("CRASH: %v", r)
			saved.ParsedOK = false
		}
		
		// Sauvegarder dans tous les cas
		data, _ := json.MarshalIndent(saved, "", "  ")
		os.WriteFile(filename, data, 0644)
		
		if saved.ParsedOK {
			fmt.Printf("✅ Parsing réussi (%d éléments) - sauvegardé dans %s\n", saved.ElementsCount, filename)
		} else {
			fmt.Printf("❌ Parsing échoué - erreur sauvegardée dans %s\n", filename)
		}
	}()

	// Tentative de parsing
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	// Si on arrive ici, pas de crash
	saved.ElementsCount = len(program)
	saved.ParsedOK = true
	
	for _, stmt := range program {
		if stmt != nil {
			saved.ParsedTypes = append(saved.ParsedTypes, fmt.Sprintf("%T", stmt))
		}
	}
}

// Lire le parsing depuis le fichier et générer les transpositions
func generateFromSavedParsing(filename string) {
	fmt.Printf("📖 Lecture du parsing depuis %s...\n", filename)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ Impossible de lire %s: %v\n", filename, err)
		return
	}

	var saved SavedParsing
	err = json.Unmarshal(data, &saved)
	if err != nil {
		fmt.Printf("❌ Impossible de décoder %s: %v\n", filename, err)
		return
	}

	fmt.Printf("📊 Parsing lu: %d éléments, Succès: %v\n", saved.ElementsCount, saved.ParsedOK)
	
	if !saved.ParsedOK {
		fmt.Printf("❌ Le parsing avait échoué: %s\n", saved.ErrorMessage)
		fmt.Println("🔧 Il faut corriger le parser avant de pouvoir générer les transpositions")
		return
	}

	if saved.ElementsCount == 0 {
		fmt.Println("⚠️ Aucun élément parsé - rien à transpiler")
		return
	}

	// Re-parser SEULEMENT si le parsing précédent était OK
	fmt.Println("🔄 Re-parsing pour génération (parsing précédent était OK)...")
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ CRASH lors de la génération: %v\n", r)
		}
	}()

	l := lexer.New(saved.SourceCode)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(program) == 0 {
		fmt.Println("❌ Re-parsing échoué")
		return
	}

	// Générer les transpositions
	languages := []generator.TargetLanguage{
		generator.JavaScript,
		generator.Java,
		generator.Python,
		generator.CSharp,
		generator.Go,
	}

	languageNames := []string{
		"JavaScript",
		"Java", 
		"Python",
		"C#",
		"Go",
	}

	// Créer fichier de transpositions
	var content strings.Builder
	content.WriteString("TRANSPOSITIONS GÉNÉRÉES\n")
	content.WriteString("========================\n\n")
	content.WriteString("CODE SOURCE:\n")
	content.WriteString(saved.SourceCode + "\n\n")

	fmt.Println("🎯 Génération des transpositions...")
	
	for i, lang := range languages {
		fmt.Printf("  📝 %s...\n", languageNames[i])
		content.WriteString(fmt.Sprintf("=== %s ===\n", languageNames[i]))
		
		output := generator.Generate(program, lang)
		if output == "" {
			content.WriteString("(Aucune sortie générée)\n")
		} else {
			content.WriteString(output)
		}
		content.WriteString("\n\n")
	}

	// Sauvegarder les transpositions
	os.WriteFile("transpositions_generees.txt", []byte(content.String()), 0644)
	fmt.Println("✅ Transpositions sauvegardées dans 'transpositions_generees.txt'")
}

func main() {
	fmt.Println("🚀 SYSTÈME DE TRANSPILATION SÉCURISÉ")
	fmt.Println("=====================================\n")

	// Code de test simple
	simpleCode := `const nom: string = "Lucie";
let age: number = 17;
var majeur: boolean = false;`

	fmt.Println("1️⃣ Test avec code simple:")
	saveParsingOnce(simpleCode, "parsing_simple.json")
	generateFromSavedParsing("parsing_simple.json")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Code complexe
	complexCode := `const nom: string = "Lucie";
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
console.log("Est majeur :", majeur);`

	fmt.Println("2️⃣ Test avec code complet:")
	saveParsingOnce(complexCode, "parsing_complex.json")
	generateFromSavedParsing("parsing_complex.json")

	fmt.Println("\n🎉 Tests terminés ! Consultez les fichiers générés :")
	fmt.Println("📄 parsing_simple.json - Résultat parsing simple")
	fmt.Println("📄 parsing_complex.json - Résultat parsing complexe") 
	fmt.Println("📄 transpositions_generees.txt - Transpositions générées")
}

package main

import (
	"ProjetGo/generator"
	"fmt"
	"regexp"
	"strings"
)

// normalizeString normalise une chaîne pour la comparaison
// en supprimant les espaces multiples et les commentaires
func normalizeString(s string) string {
	// Normaliser les sauts de ligne
	s = strings.ReplaceAll(s, "\r\n", "\n")
	
	// Traiter ligne par ligne
	lines := strings.Split(s, "\n")
	normalizedLines := make([]string, 0, len(lines))
	
	for _, line := range lines {
		// Supprimer les espaces en début et fin de ligne
		line = strings.TrimSpace(line)
		
		// Ignorer les lignes vides
		if line == "" {
			continue
		}
		
		// Normaliser les espaces multiples
		spaceRegex := regexp.MustCompile(`\s+`)
		line = spaceRegex.ReplaceAllString(line, " ")
		
		// Ajouter la ligne normalisée
		normalizedLines = append(normalizedLines, line)
	}
	
	// Rejoindre les lignes
	return strings.Join(normalizedLines, "\n")
}

// TestTranspilation exécute un test de transpilation TypeScript → JavaScript
func TestTranspilation() {
	fmt.Println("TEST: Fonction avec paramètre par défaut")
	fmt.Println("========================================")
	
	input := `
function incrementer(nombre: number, pas: number = 1): number {
  return nombre + pas;
}

// Exemple d'utilisation
console.log(incrementer(5)); // Devrait afficher 6
console.log(incrementer(5, 2)); // Devrait afficher 7
`

	expected := `
function incrementer(nombre, pas = 1) {
  return nombre + pas;
}

// Exemple d'utilisation
console.log(incrementer(5)); // Devrait afficher 6
console.log(incrementer(5, 2)); // Devrait afficher 7
`

	// Transpilation avec la méthode automatique
	autoOutput := generator.TranspileTS(input)
	
	fmt.Println("\nTypeScript original:")
	fmt.Println("-------------------")
	fmt.Println(input)
	
	fmt.Println("\nJavaScript généré (automatique):")
	fmt.Println("-----------------------------")
	fmt.Println(autoOutput)
	
	fmt.Println("\nJavaScript attendu:")
	fmt.Println("------------------")
	fmt.Println(expected)
	
	// Normaliser les chaînes pour la comparaison
	normalizedAuto := normalizeString(autoOutput)
	normalizedExpected := normalizeString(expected)
	
	// Afficher si la sortie automatique correspond à l'attendu
	fmt.Println("\nRÉSULTATS:")
	fmt.Println("=========")
	
	if normalizedAuto == normalizedExpected {
		fmt.Println("✅ La transpilation est fonctionnelle!")
	} else {
		fmt.Println("❌ La transpilation a des problèmes")
	}
	
	// Test 2: Classe
	fmt.Println("\n\nTEST: Classe avec méthode ayant des paramètres par défaut")
	fmt.Println("======================================================")
	
	input2 := `
class Calculateur {
  private valeur: number;
  
  constructor(valeurInitiale: number = 0) {
    this.valeur = valeurInitiale;
  }
  
  incrementer(pas: number = 1): void {
    this.valeur += pas;
  }
  
  getValeur(): number {
    return this.valeur;
  }
}

// Exemple d'utilisation
const calc = new Calculateur(10);
calc.incrementer(); // +1
calc.incrementer(5); // +5
console.log(calc.getValeur()); // Devrait afficher 16
`

	expected2 := `
class Calculateur {
  constructor(valeurInitiale = 0) {
    this.valeur = valeurInitiale;
  }
  
  incrementer(pas = 1) {
    this.valeur += pas;
  }
  
  getValeur() {
    return this.valeur;
  }
}

// Exemple d'utilisation
const calc = new Calculateur(10);
calc.incrementer(); // +1
calc.incrementer(5); // +5
console.log(calc.getValeur()); // Devrait afficher 16
`

	// Transpilation avec la méthode automatique
	autoOutput2 := generator.TranspileTS(input2)
	
	fmt.Println("\nTypeScript original:")
	fmt.Println("-------------------")
	fmt.Println(input2)
	
	fmt.Println("\nJavaScript généré (automatique):")
	fmt.Println("-----------------------------")
	fmt.Println(autoOutput2)
	
	fmt.Println("\nJavaScript attendu:")
	fmt.Println("------------------")
	fmt.Println(expected2)
	
	// Normaliser les chaînes pour la comparaison
	normalizedAuto2 := normalizeString(autoOutput2)
	normalizedExpected2 := normalizeString(expected2)
	
	if normalizedAuto2 == normalizedExpected2 {
		fmt.Println("\n✅ La transpilation des classes est fonctionnelle!")
	} else {
		fmt.Println("\n❌ La transpilation des classes a des problèmes")
	}
}

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

func main() {
	// Test 1: Fonction avec paramètre par défaut
	fmt.Println("TEST 1: Fonction avec paramètre par défaut")
	fmt.Println("==========================================")
	
	input1 := `
function incrementer(nombre: number, pas: number = 1): number {
  return nombre + pas;
}

// Exemple d'utilisation
console.log(incrementer(5)); // Devrait afficher 6
console.log(incrementer(5, 2)); // Devrait afficher 7
`

	expected1 := `
function incrementer(nombre, pas = 1) {
  return nombre + pas;
}

// Exemple d'utilisation
console.log(incrementer(5)); // Devrait afficher 6
console.log(incrementer(5, 2)); // Devrait afficher 7
`

	autoOutput1 := generator.TranspileTS(input1)
	
	fmt.Println("\nTypeScript original:")
	fmt.Println("-------------------")
	fmt.Println(input1)
	
	fmt.Println("\nJavaScript généré (automatique):")
	fmt.Println("-----------------------------")
	fmt.Println(autoOutput1)
	
	fmt.Println("\nJavaScript attendu:")
	fmt.Println("------------------")
	fmt.Println(expected1)
	
	// Test 2: Classe avec méthode ayant des paramètres par défaut
	fmt.Println("\n\nTEST 2: Classe avec méthode ayant des paramètres par défaut")
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
	
	// Comparer les résultats normalisés
	fmt.Println("\n\nRÉSULTATS:")
	fmt.Println("=========")
	
	normalizedAuto1 := normalizeString(autoOutput1)
	normalizedExpected1 := normalizeString(expected1)
	
	if normalizedAuto1 == normalizedExpected1 {
		fmt.Println("✅ Test 1: La sortie correspond à l'attendu!")
	} else {
		fmt.Println("❌ Test 1: La sortie diffère de l'attendu")
	}
	
	normalizedAuto2 := normalizeString(autoOutput2)
	normalizedExpected2 := normalizeString(expected2)
	
	if normalizedAuto2 == normalizedExpected2 {
		fmt.Println("✅ Test 2: La sortie correspond à l'attendu!")
	} else {
		fmt.Println("❌ Test 2: La sortie diffère de l'attendu")
	}
}

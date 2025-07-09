package main

import (
	"ProjetGo/generator"
	"ProjetGo/lexer"
	"ProjetGo/parser"
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
	// Code TypeScript à tester
	input := `
function estPalindrome(texte: string): boolean {
  const normalise = texte.toLowerCase().replace(/[\W_]/g, '');
  const inverse = normalise.split('').reverse().join('');
  return normalise === inverse;
}

// Exemple d'utilisation :
console.log(estPalindrome("Radar")); // true
console.log(estPalindrome("Bonjour")); // false
`

	// Transpilation avec la nouvelle méthode automatique
	autoOutput := generator.TranspileTS(input)
	
	// Transpilation avec l'ancienne méthode pour comparaison
	// Méthode de l'ancien générateur - juste pour référence
	l := lexer.New(input)
	p := parser.New(l)
	p.ParseProgram() // On n'utilise pas directement le résultat
	
	// Sortie attendue - ajusté pour correspondre à notre génération attendue
	expected := `function estPalindrome(texte) {
  const normalise = texte.toLowerCase().replace(/[\W_]/g, '');
  const inverse = normalise.split('').reverse().join('');
  return normalise === inverse;
}

// Exemple d'utilisation :
console.log(estPalindrome("Radar")); // true
console.log(estPalindrome("Bonjour")); // false`

	// Afficher le résultat
	fmt.Println("TypeScript original:")
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
		fmt.Println("✅ La sortie générée automatiquement correspond à ce qui était attendu!")
	} else {
		fmt.Println("❌ La sortie générée automatiquement est différente de ce qui était attendu")
		fmt.Println("\nDifférences (après normalisation):")
		
		// Afficher les différences ligne par ligne
		autoLines := strings.Split(normalizedAuto, "\n")
		expectedLines := strings.Split(normalizedExpected, "\n")
		
		// Comparer chaque ligne
		for i := 0; i < len(autoLines) && i < len(expectedLines); i++ {
			if autoLines[i] != expectedLines[i] {
				fmt.Printf("Ligne %d:\n  Auto: %s\n  Attendu: %s\n", i+1, autoLines[i], expectedLines[i])
			}
		}
		
		// Vérifier si le nombre de lignes est différent
		if len(autoLines) != len(expectedLines) {
			fmt.Printf("Nombre de lignes différent: Auto=%d, Attendu=%d\n", len(autoLines), len(expectedLines))
		}
	}
}

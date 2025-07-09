package main

// 🚀 TRANSPILEUR TYPESCRIPT → JAVASCRIPT - GO PURE
// ================================================
// 
// ✨ ZERO dépendance externe - 100% bibliothèque standard Go
//
// Imports utilisés (tous de la stdlib Go) :
// - ProjetGo/generator : Notre générateur de code (interne)
// - ProjetGo/lexer     : Notre analyseur lexical (interne)  
// - ProjetGo/parser    : Notre analyseur syntaxique (interne)
// - fmt                : Formatage et affichage
// - log                : Logging
// - net/http           : Serveur web HTTP
// - os                 : Opérations système
// - encoding/json      : Parsing/génération JSON (utilisé plus bas)
// - strings            : Manipulation de chaînes (utilisé plus bas)
// - html/template      : Templates HTML (utilisé plus bas)

import (
	"ProjetGo/generator"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"flag"
	"fmt"
	"os"
	"time"
)

// Configuration structure
type Config struct {
	Port    string
	Host    string
	Console bool
	File    string
	Target  string
	Verbose bool
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Port, "port", "8080", "Port for web server")
	flag.StringVar(&config.Host, "host", "localhost", "Host for web server")
	flag.BoolVar(&config.Console, "console", false, "Run in console mode")
	flag.StringVar(&config.File, "file", "", "Input file to transpile")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose output")

	flag.Parse()

	// Handle legacy console argument
	if len(os.Args) > 1 && os.Args[1] == "console" {
		config.Console = true
	}

	return config
}

func runConsoleVersion(config *Config) {
	var input string

	if config.File != "" {
		// Read from file
		content, err := os.ReadFile(config.File)
		if err != nil {
			fmt.Printf("❌ Erreur lors de la lecture du fichier: %v\n", err)
			return
		}
		input = string(content)
	} else {
		// Default example
		input = `
const message: string = "Hello World";
let count: number = 42;
const pi: number = 3.14;
let isActive: boolean = true;

function greet(name: string): string {
    return "Hello " + name;
}

interface User {
    id: number;
    name: string;
    email: string;
}

class Calculator {
    private value: number = 0;
    
    add(x: number): void {
        this.value += x;
    }
    
    getResult(): number {
        return this.value;
    }
}
`
	}

	fmt.Println("🚀 Transpilateur TypeScript → JavaScript - Mode Console")
	fmt.Println("==================================================")
	fmt.Println("📝 Code Source (TypeScript-like):")
	fmt.Println("--------------------------------")
	fmt.Println(input)
	fmt.Println()

	start := time.Now()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	if len(program.Statements) == 0 {
		fmt.Println("❌ Aucun code valide détecté. Vérifiez la syntaxe.")
		return
	}

	elapsed := time.Since(start)
	if config.Verbose {
		fmt.Printf("⏱️  Temps de parsing: %v\n\n", elapsed)
	}

	// Générer le code JavaScript à partir du code source directement
	fmt.Println("=== 🟨 JavaScript Output ===")
	fmt.Println(generator.TranspileTS(input))
	fmt.Println()

	totalElapsed := time.Since(start)
	if config.Verbose {
		fmt.Printf("⏱️  Temps total: %v\n", totalElapsed)
	}
}

// Ces fonctions ne sont plus nécessaires car nous ne gérons que TypeScript vers JavaScript

func main() {
	config := parseFlags()

	fmt.Println("🚀 Transpilateur TypeScript → JavaScript v2.0")
	fmt.Println("============================================")

	// Option pour exécuter le test de transpilation
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("🧪 Exécution du test de transpilation...")
		fmt.Println("Pour exécuter les tests, utilisez 'go run test_default_params.go' ou 'go test'")
		return
	}

	if config.Console {
		fmt.Println("💻 Mode console activé...")
		runConsoleVersion(config)
		return
	}

	// Par défaut, lancer l'interface web
	fmt.Println("🌐 Lancement de l'interface web...")
	fmt.Println("💡 Pour utiliser la version console: go run . console")
	fmt.Println("💡 Pour exécuter le test de transpilation: go run . test")
	fmt.Println("")
	
	StartWebServer()
}

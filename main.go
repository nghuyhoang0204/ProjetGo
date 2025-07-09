package main

import (
	"ProjetGo/generator"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"flag"
	"fmt"
	"os"
	"strings"
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
	flag.StringVar(&config.Target, "target", "all", "Target language (js,java,python,csharp,go,rust,swift,php,all)")
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

	fmt.Println("🚀 Transpilateur Multi-Langages - Mode Console")
	fmt.Println("=============================================")
	fmt.Println("📝 Code Source (TypeScript-like):")
	fmt.Println("--------------------------------")
	fmt.Println(input)
	fmt.Println()

	start := time.Now()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	if len(program) == 0 {
		fmt.Println("❌ Aucun code valide détecté. Vérifiez la syntaxe.")
		return
	}

	elapsed := time.Since(start)
	if config.Verbose {
		fmt.Printf("⏱️  Temps de parsing: %v\n\n", elapsed)
	}

	// Generate code for specified targets
	targets := getTargetLanguages(config.Target)

	for _, target := range targets {
		fmt.Printf("=== %s Output ===\n", getLanguageName(target))
		fmt.Println(generator.Generate(program, target))
		fmt.Println()
	}

	totalElapsed := time.Since(start)
	if config.Verbose {
		fmt.Printf("⏱️  Temps total: %v\n", totalElapsed)
	}
}

func getTargetLanguages(target string) []generator.TargetLanguage {
	switch strings.ToLower(target) {
	case "js", "javascript":
		return []generator.TargetLanguage{generator.JavaScript}
	case "java":
		return []generator.TargetLanguage{generator.Java}
	case "python":
		return []generator.TargetLanguage{generator.Python}
	case "csharp", "c#":
		return []generator.TargetLanguage{generator.CSharp}
	case "go":
		return []generator.TargetLanguage{generator.Go}
	case "rust":
		return []generator.TargetLanguage{generator.Rust}
	case "swift":
		return []generator.TargetLanguage{generator.Swift}
	case "php":
		return []generator.TargetLanguage{generator.PHP}
	default:
		return []generator.TargetLanguage{
			generator.JavaScript,
			generator.Java,
			generator.Python,
			generator.CSharp,
			generator.Go,
			generator.Rust,
			generator.Swift,
			generator.PHP,
		}
	}
}

func getLanguageName(target generator.TargetLanguage) string {
	switch target {
	case generator.JavaScript:
		return "🟨 JavaScript"
	case generator.Java:
		return "☕ Java"
	case generator.Python:
		return "🐍 Python"
	case generator.CSharp:
		return "🔵 C#"
	case generator.Go:
		return "🐹 Go"
	case generator.Rust:
		return "🦀 Rust"
	case generator.Swift:
		return "🍎 Swift"
	case generator.PHP:
		return "🐘 PHP"
	default:
		return "Unknown"
	}
}

func main() {
	config := parseFlags()

	fmt.Println("🚀 Transpilateur Multi-Langages v2.0")
	fmt.Println("====================================")

	if config.Console {
		fmt.Println("💻 Mode console activé...")
		runConsoleVersion(config)
		return
	}

	// Set server configuration
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	fmt.Println("🌐 Lancement de l'interface web...")
	fmt.Printf("📍 Serveur: http://%s\n", serverAddr)
	fmt.Println("💡 Commandes disponibles:")
	fmt.Println("   go run . --console          # Mode console")
	fmt.Println("   go run . --file=code.ts     # Transpiler un fichier")
	fmt.Println("   go run . --target=python    # Langage cible spécifique")
	fmt.Println("   go run . --port=3000        # Port personnalisé")
	fmt.Println("   go run . --verbose          # Sortie détaillée")
	fmt.Println()

	StartWebServer()
}

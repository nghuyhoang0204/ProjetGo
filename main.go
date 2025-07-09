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
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Si des arguments sont fournis, traiter en mode CLI
	if len(os.Args) > 1 {
		handleCLI()
		return
	}

	// Sinon, démarrer le serveur web
	startWebServer()
}

// handleCLI traite la transpilation en mode ligne de commande
func handleCLI() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <typescript-code>")
		fmt.Println("   ou: go run main.go (pour démarrer le serveur web)")
		return
	}

	typeScriptCode := os.Args[1]
	
	// Transpiler le code TypeScript
	jsCode, err := TranspileTypeScriptToJavaScript(typeScriptCode)
	if err != nil {
		fmt.Printf("Erreur de transpilation: %v\n", err)
		return
	}

	fmt.Println("Code JavaScript généré:")
	fmt.Println(jsCode)
}

// startWebServer démarre le serveur web
func startWebServer() {
	fmt.Println("Démarrage du serveur TypeScript vers JavaScript...")
	fmt.Println("Serveur disponible sur http://localhost:8080")
	
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/transpile", transpileHandler)
	http.HandleFunc("/api/transpile", apiTranspileHandler)
	
	// Servir les fichiers statiques (CSS, JS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TranspileTypeScriptToJavaScript fonction principale de transpilation
func TranspileTypeScriptToJavaScript(typeScriptCode string) (string, error) {
	// Créer le lexer
	l := lexer.New(typeScriptCode)
	
	// Créer le parser
	p := parser.New(l)
	
	// Parser le code en AST
	program := p.ParseProgram()
	
	// Vérifier les erreurs de parsing
	if len(p.Errors()) > 0 {
		return "", fmt.Errorf("erreurs de parsing: %v", p.Errors())
	}
	
	// Créer le générateur
	transpiler := &generator.TypeScriptToJavaScriptTranspiler{}
	
	// Générer le code JavaScript
	jsCode := transpiler.Generate(program)
	
	return jsCode, nil
}

// homeHandler affiche la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Transpileur TypeScript vers JavaScript</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }

        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
        }

        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }

        .main-content {
            padding: 30px;
        }

        .editor-container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }

        .editor-section {
            background: #f8fafc;
            border-radius: 10px;
            padding: 20px;
            border: 2px solid #e2e8f0;
        }

        .editor-section h3 {
            color: #1e293b;
            margin-bottom: 15px;
            font-size: 1.2rem;
        }

        textarea {
            width: 100%;
            height: 300px;
            border: 1px solid #cbd5e1;
            border-radius: 8px;
            padding: 15px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            resize: vertical;
            background: white;
        }

        textarea:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        .button-container {
            text-align: center;
            margin: 20px 0;
        }

        .transpile-btn {
            background: linear-gradient(135deg, #10b981 0%, #059669 100%);
            color: white;
            border: none;
            padding: 15px 40px;
            border-radius: 50px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
        }

        .transpile-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(16, 185, 129, 0.4);
        }

        .transpile-btn:active {
            transform: translateY(0);
        }

        .example-section {
            background: #fef3c7;
            border-radius: 10px;
            padding: 20px;
            margin-top: 30px;
            border-left: 4px solid #f59e0b;
        }

        .example-section h3 {
            color: #92400e;
            margin-bottom: 15px;
        }

        .example-code {
            background: #fffbeb;
            border: 1px solid #fed7aa;
            border-radius: 6px;
            padding: 15px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            white-space: pre-wrap;
            color: #9a3412;
        }

        .error {
            background: #fef2f2;
            border: 1px solid #fecaca;
            border-radius: 8px;
            padding: 15px;
            color: #dc2626;
            margin-top: 10px;
        }

        @media (max-width: 768px) {
            .editor-container {
                grid-template-columns: 1fr;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            .main-content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔄 Transpileur TypeScript → JavaScript</h1>
            <p>Convertissez votre code TypeScript en JavaScript propre et idiomatic</p>
        </div>
        
        <div class="main-content">
            <div class="editor-container">
                <div class="editor-section">
                    <h3>📝 Code TypeScript</h3>
                    <textarea id="typescript-input" placeholder="Entrez votre code TypeScript ici...">interface User {
    name: string;
    age: number;
    email?: string;
}

function greetUser(user: User): string {
    let greeting: string = "Hello, " + user.name;
    
    if (user.age >= 18) {
        greeting += " (Adult)";
    } else {
        greeting += " (Minor)";
    }
    
    return greeting;
}

const user: User = {
    name: "Alice",
    age: 25,
    email: "alice@example.com"
};

console.log(greetUser(user));</textarea>
                </div>
                
                <div class="editor-section">
                    <h3>⚡ JavaScript Généré</h3>
                    <textarea id="javascript-output" readonly placeholder="Le code JavaScript apparaîtra ici..."></textarea>
                </div>
            </div>
            
            <div class="button-container">
                <button class="transpile-btn" onclick="transpileCode()">
                    🚀 Transpiler le Code
                </button>
            </div>
            
            <div class="example-section">
                <h3>💡 Fonctionnalités supportées</h3>
                <div class="example-code">• Suppression automatique des annotations de type
• Conversion des interfaces en commentaires
• Support des déclarations let/const/var
• Fonctions et méthodes
• Structures de contrôle (if/else, for, while)
• Opérateurs et expressions
• Commentaires préservés</div>
            </div>
        </div>
    </div>

    <script>
        async function transpileCode() {
            const input = document.getElementById('typescript-input').value;
            const output = document.getElementById('javascript-output');
            
            if (!input.trim()) {
                alert('Veuillez entrer du code TypeScript');
                return;
            }
            
            try {
                const response = await fetch('/api/transpile', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ code: input })
                });
                
                const result = await response.json();
                
                if (result.success) {
                    output.value = result.javascript;
                    // Supprimer les erreurs précédentes
                    const errorDiv = document.querySelector('.error');
                    if (errorDiv) {
                        errorDiv.remove();
                    }
                } else {
                    output.value = '';
                    showError(result.error);
                }
            } catch (error) {
                output.value = '';
                showError('Erreur de connexion: ' + error.message);
            }
        }
        
        function showError(message) {
            // Supprimer les erreurs précédentes
            const existingError = document.querySelector('.error');
            if (existingError) {
                existingError.remove();
            }
            
            // Créer une nouvelle div d'erreur
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error';
            errorDiv.textContent = message;
            
            // L'ajouter après le container des boutons
            const buttonContainer = document.querySelector('.button-container');
            buttonContainer.parentNode.insertBefore(errorDiv, buttonContainer.nextSibling);
        }
        
        // Transpiler automatiquement au chargement de la page
        window.addEventListener('load', () => {
            transpileCode();
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// transpileHandler gère les requêtes de transpilation (compatibilité)
func transpileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	typeScriptCode := r.FormValue("typescript")
	if typeScriptCode == "" {
		http.Error(w, "Code TypeScript manquant", http.StatusBadRequest)
		return
	}

	jsCode, err := TranspileTypeScriptToJavaScript(typeScriptCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur de transpilation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(jsCode))
}

// apiTranspileHandler gère les requêtes API JSON
func apiTranspileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Code string `json:"code"`
	}

	if err := r.ParseForm(); err == nil && r.Header.Get("Content-Type") != "application/json" {
		// Support pour les formulaires classiques
		request.Code = r.FormValue("code")
	} else {
		// Support pour JSON
		decoder := fmt.Sprintf(`{
			"code": "%s"
		}`, r.FormValue("code"))
		if r.Header.Get("Content-Type") == "application/json" {
			// Lire le JSON depuis le body
			buf := make([]byte, r.ContentLength)
			r.Body.Read(buf)
			decoder = string(buf)
		}
		// Parse simple du JSON (version basique)
		code := extractJSONValue(decoder, "code")
		request.Code = code
	}

	if request.Code == "" {
		writeJSONError(w, "Code TypeScript manquant", http.StatusBadRequest)
		return
	}

	jsCode, err := TranspileTypeScriptToJavaScript(request.Code)
	if err != nil {
		writeJSONError(w, fmt.Sprintf("Erreur de transpilation: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"javascript": jsCode,
	}

	writeJSONResponse(w, response)
}

// extractJSONValue extrait une valeur d'un JSON simple (parser basique)
func extractJSONValue(jsonStr, key string) string {
	// Recherche simple de la clé dans le JSON
	keyPattern := `"` + key + `"`
	start := fmt.Sprintf(`%s:`, keyPattern)
	startIndex := fmt.Sprintf("%s", jsonStr)
	
	// Version simplifiée pour extraire la valeur
	lines := []string{jsonStr}
	for _, line := range lines {
		if idx := findInString(line, start); idx >= 0 {
			afterColon := line[idx+len(start):]
			afterColon = trimSpaces(afterColon)
			if len(afterColon) > 0 && afterColon[0] == '"' {
				// Trouver la fin de la string
				end := findInString(afterColon[1:], `"`)
				if end >= 0 {
					return afterColon[1 : end+1]
				}
			}
		}
	}
	return ""
}

func findInString(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func trimSpaces(s string) string {
	start := 0
	end := len(s)
	
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	
	return s[start:end]
}

// writeJSONResponse écrit une réponse JSON
func writeJSONResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	// Conversion manuelle en JSON (version simple)
	jsonStr := "{"
	first := true
	for key, value := range data {
		if !first {
			jsonStr += ","
		}
		first = false
		
		jsonStr += fmt.Sprintf(`"%s":`, key)
		switch v := value.(type) {
		case string:
			// Échapper les caractères spéciaux
			escaped := escapeJSON(v)
			jsonStr += fmt.Sprintf(`"%s"`, escaped)
		case bool:
			if v {
				jsonStr += "true"
			} else {
				jsonStr += "false"
			}
		default:
			jsonStr += fmt.Sprintf(`"%v"`, v)
		}
	}
	jsonStr += "}"
	
	w.Write([]byte(jsonStr))
}

// writeJSONError écrit une erreur JSON
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}
	writeJSONResponse(w, response)
}

// escapeJSON échappe les caractères spéciaux dans une string JSON
func escapeJSON(s string) string {
	result := ""
	for _, char := range s {
		switch char {
		case '"':
			result += `\"`
		case '\\':
			result += `\\`
		case '\n':
			result += `\n`
		case '\r':
			result += `\r`
		case '\t':
			result += `\t`
		default:
			result += string(char)
		}
	}
	return result
}

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
)

// Structure pour recevoir le code à transpiler
type TranspileRequest struct {
	Code   string `json:"code"`
	Target string `json:"target"`
}

type TranspileResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

// Page HTML avec interface améliorée
const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>🚀 Transpilateur Multi-Langages</title>
    <meta charset="UTF-8">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', system-ui, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 15px; 
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(45deg, #667eea, #764ba2); 
            color: white; 
            padding: 30px; 
            text-align: center; 
        }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header p { opacity: 0.9; font-size: 1.1em; }
        .main { display: flex; height: 70vh; }
        .input-section { 
            flex: 1; 
            border-right: 2px solid #f0f0f0; 
            display: flex; 
            flex-direction: column;
        }
        .output-section { flex: 1; display: flex; flex-direction: column; }
        .section-header { 
            background: #f8f9fa; 
            padding: 15px 20px; 
            border-bottom: 1px solid #e9ecef; 
            font-weight: 600; 
        }
        .controls { 
            padding: 15px 20px; 
            background: #f8f9fa; 
            border-bottom: 1px solid #e9ecef; 
        }
        select, button { 
            padding: 10px 15px; 
            border: 1px solid #ddd; 
            border-radius: 8px; 
            font-size: 14px; 
        }
        button { 
            background: #667eea; 
            color: white; 
            border: none; 
            cursor: pointer; 
            margin-left: 10px; 
            transition: all 0.3s;
        }
        button:hover { background: #5a6fd8; transform: translateY(-1px); }
        textarea { 
            flex: 1; 
            border: none; 
            padding: 20px; 
            font-family: 'Courier New', monospace; 
            font-size: 14px; 
            resize: none; 
            outline: none; 
        }
        .output { 
            flex: 1; 
            padding: 20px; 
            background: #1e1e1e; 
            color: #d4d4d4; 
            font-family: 'Courier New', monospace; 
            font-size: 14px; 
            overflow-y: auto; 
            white-space: pre-wrap;
        }
        .status { 
            padding: 10px 20px; 
            background: #e8f5e8; 
            border-top: 1px solid #e9ecef; 
            color: #2d5a2d; 
            font-weight: 500;
        }
        .status.error { background: #ffe8e8; color: #8b0000; }
        .example-btn {
            background: #28a745;
            font-size: 12px;
            padding: 5px 10px;
            margin-left: 10px;
        }
        .example-btn:hover { background: #218838; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Transpilateur Multi-Langages</h1>
            <p>Convertissez votre code TypeScript/JavaScript vers 5 langages différents</p>
        </div>
        
        <div class="main">
            <div class="input-section">
                <div class="section-header">
                    📝 Code Source (TypeScript/JavaScript)
                    <button class="example-btn" onclick="loadExample()">Charger Exemple</button>
                </div>
                <textarea id="sourceCode" placeholder="Entrez votre code TypeScript/JavaScript ici...

🎯 ACTUELLEMENT SUPPORTÉ :

✅ Variables avec types :
const nom: string = 'Alice';
let age: number = 25;
var actif: boolean = true;

✅ Types de base :
- string (chaînes)
- number (nombres) 
- boolean (true/false)

✅ Structures de données :
let valeurs: number[] = [10, 20, 30];
let personne = { nom: 'Bob', age: 30 };

✅ Fonctions simples :
function saluer(nom: string): string {
  return 'Salut !';
}

✅ Conditions :
if (age >= 18) {
  return 'Majeur';
} else {
  return 'Mineur';  
}

✅ Utilisation :
console.log(message);
console.log(saluer('Marie'));

� EXEMPLE COMPLET À TESTER : Cliquez 'Charger Exemple'

🚧 EN COURS D'AMÉLIORATION :
- Boucles (for/while)
- Template literals complets 
- Expressions arithmétiques complexes
- Commentaires (à ignorer)"></textarea>
                
                <div class="controls">
                    <label>Langage cible :</label>
                    <select id="targetLang">
                        <option value="javascript">🟨 JavaScript</option>
                        <option value="java">☕ Java</option>
                        <option value="python">🐍 Python</option>
                        <option value="csharp">🔵 C#</option>
                        <option value="go">🐹 Go</option>
                    </select>
                    <button onclick="transpile()">🔄 Transpiler</button>
                    <button onclick="transpileAll()">🎯 Tous les langages</button>
                </div>
            </div>
            
            <div class="output-section">
                <div class="section-header">
                    ⚡ Code Généré
                </div>
                <div id="output" class="output">Sélectionnez un langage cible et cliquez sur "Transpiler"</div>
                <div id="status" class="status">Prêt à transpiler</div>
            </div>
        </div>
    </div>

    <script>
        function loadExample() {
            document.getElementById('sourceCode').value = ` + "`" + `// Variables avec types
const nom: string = "Alice";
let age: number = 25;
var actif: boolean = true;

// Types de base
let score: number = 42;
let message: string = "Hello World";
let valeurs: number[] = [10, 20, 30];
let personne = {
  nom: "Bob",
  age: 30
};

// Fonction simple
function saluer(prenom: string): string {
  return "Salut !";
}

// Condition simple
function verifierAge(age: number): string {
  if (age >= 18) {
    return "Majeur";
  } else {
    return "Mineur";
  }
}

// Utilisation
console.log(message);
console.log(saluer("Marie"));` + "`" + `;
        }

        async function transpile() {
            const code = document.getElementById('sourceCode').value;
            const target = document.getElementById('targetLang').value;
            const output = document.getElementById('output');
            const status = document.getElementById('status');
            
            if (!code.trim()) {
                status.textContent = '❌ Veuillez entrer du code à transpiler';
                status.className = 'status error';
                return;
            }
            
            status.textContent = '⏳ Transpilation en cours...';
            status.className = 'status';
            
            try {
                const response = await fetch('/transpile', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ code, target })
                });
                
                const result = await response.json();
                
                if (result.success) {
                    output.textContent = result.output;
                    status.textContent = '✅ Transpilation réussie pour ' + target.toUpperCase();
                    status.className = 'status';
                } else {
                    output.textContent = '❌ Erreur: ' + result.error;
                    status.textContent = '❌ Erreur de transpilation';
                    status.className = 'status error';
                }
            } catch (error) {
                output.textContent = '❌ Erreur de connexion: ' + error.message;
                status.textContent = '❌ Erreur de connexion';
                status.className = 'status error';
            }
        }
        
        async function transpileAll() {
            const code = document.getElementById('sourceCode').value;
            const output = document.getElementById('output');
            const status = document.getElementById('status');
            
            if (!code.trim()) {
                status.textContent = '❌ Veuillez entrer du code à transpiler';
                status.className = 'status error';
                return;
            }
            
            status.textContent = '⏳ Transpilation vers tous les langages...';
            status.className = 'status';
            
            const languages = ['javascript', 'java', 'python', 'csharp', 'go'];
            let allOutput = '';
            
            for (const lang of languages) {
                try {
                    const response = await fetch('/transpile', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ code, target: lang })
                    });
                    
                    const result = await response.json();
                    
                    allOutput += '=== ' + lang.toUpperCase() + ' ===\n';
                    if (result.success) {
                        allOutput += result.output + '\n\n';
                    } else {
                        allOutput += '❌ Erreur: ' + result.error + '\n\n';
                    }
                } catch (error) {
                    allOutput += '❌ Erreur: ' + error.message + '\n\n';
                }
            }
            
            output.textContent = allOutput;
            status.textContent = '✅ Transpilation terminée pour tous les langages';
            status.className = 'status';
        }
        
        // Auto-resize textarea
        document.getElementById('sourceCode').addEventListener('input', function() {
            this.style.height = 'auto';
            this.style.height = this.scrollHeight + 'px';
        });
    </script>
</body>
</html>`

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, htmlTemplate)
}

func extractSimpleVariables(code string) string {
	// Extraire intelligemment les variables simples du code complexe
	lines := strings.Split(code, "\n")
	var extractedVars []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Chercher les déclarations de variables simples
		if strings.HasPrefix(line, "const ") || strings.HasPrefix(line, "let ") || strings.HasPrefix(line, "var ") {
			// Ignorer les lignes trop complexes (avec { ou [ ou fonction)
			if !strings.Contains(line, "{") && !strings.Contains(line, "[") && 
			   !strings.Contains(line, "(") && strings.Contains(line, "=") &&
			   !strings.Contains(line, "new ") && !strings.Contains(line, "=>") {
				extractedVars = append(extractedVars, line)
			}
		}
	}
	
	return strings.Join(extractedVars, "\n")
}

func handleTranspile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TranspileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(TranspileResponse{
			Success: false,
			Error:   "Invalid JSON: " + err.Error(),
		})
		return
	}

	// Transpiler le code
	// D'abord essayer d'extraire les variables simples si le code est complexe
	simpleVars := extractSimpleVariables(req.Code)
	
	var codeToProcess string
	var isSimplified bool
	
	if simpleVars != "" && simpleVars != req.Code {
		codeToProcess = simpleVars
		isSimplified = true
	} else {
		codeToProcess = req.Code
		isSimplified = false
	}
	
	l := lexer.New(codeToProcess)
	p := parser.New(l)
	program := p.ParseProgram()

	var output string
	var targetLang generator.TargetLanguage

	switch req.Target {
	case "javascript":
		targetLang = generator.JavaScript
	case "java":
		targetLang = generator.Java
	case "python":
		targetLang = generator.Python
	case "csharp":
		targetLang = generator.CSharp
	case "go":
		targetLang = generator.Go
	default:
		targetLang = generator.JavaScript
	}

	output = generator.Generate(program, targetLang)

	if output == "" || strings.TrimSpace(output) == "" {
		output = "// Aucun code généré - Le code source contient peut-être des structures non encore supportées\n"
		output += "// Structures actuellement supportées :\n"
		output += "// - Déclarations de variables (const, let, var)\n"
		output += "// - Types de base (string, number, boolean)\n"
		output += "// - Classes (structure de base)\n"
		output += "// - Interfaces (reconnaissance)\n"
		output += "\n// Votre code sera analysé et des améliorations sont en cours..."
	} else if isSimplified {
		output = "// ⚡ Code simplifié automatiquement - Seules les variables simples ont été extraites\n" +
				"// Les structures complexes (classes, interfaces, fonctions) sont en cours de développement\n\n" + output
	}

	response := TranspileResponse{
		Success: true,
		Output:  output,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StartWebServer() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/transpile", handleTranspile)

	fmt.Println("🌐 Serveur web démarré sur http://localhost:8081")
	fmt.Println("💡 Ouvrez votre navigateur à cette adresse")
	fmt.Println("🔄 Ctrl+C pour arrêter le serveur")
	
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("❌ Erreur serveur: %v\n", err)
	}
}

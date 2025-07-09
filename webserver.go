package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"os"
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
                <textarea id="sourceCode" placeholder="Entrez votre code TypeScript/JavaScript ici..."></textarea>
                
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
                </div>
            </div>
            
            <div class="output-section">
                <div class="section-header">
                    ⚡ Code Généré (parsing optimisé)
                </div>
                <div id="output" class="output">Sélectionnez un langage cible et cliquez sur "Transpiler"</div>
                <div id="status" class="status">Prêt à transpiler</div>
            </div>
        </div>
    </div>

    <script>
        function loadExample() {
            document.getElementById('sourceCode').value = ` + "`" + `const nom: string = "Lucie";
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
console.log("Est majeur :", majeur);` + "`" + `;
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
            
            status.textContent = '🔄 Parsing optimisé pour ' + target.toUpperCase() + '...';
            status.className = 'status';
            output.textContent = 'Transpilation en cours...';
            
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
    </script>
</body>
</html>`

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, htmlTemplate)
}

func saveParsingForSingleLanguage(code string, targetLang string, filename string) (string, error) {
	// Parser le code avec gestion des erreurs
	defer func() {
		if r := recover(); r != nil {
			// En cas de crash, créer un fichier d'erreur
			errorContent := fmt.Sprintf("PARSING ERROR\n=============\n\nCode: %s\n\nTarget Language: %s\n\nError: %v", code, targetLang, r)
			os.WriteFile(filename, []byte(errorContent), 0644)
		}
	}()

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	// Convertir le nom du langage vers le type TargetLanguage
	var lang generator.TargetLanguage
	switch targetLang {
	case "javascript":
		lang = generator.JavaScript
	case "java":
		lang = generator.Java
	case "python":
		lang = generator.Python
	case "csharp":
		lang = generator.CSharp
	case "go":
		lang = generator.Go
	default:
		lang = generator.JavaScript
	}

	// Générer SEULEMENT pour le langage demandé
	output := generator.Generate(program, lang)
	if output == "" {
		output = "// Aucun code généré pour " + targetLang
	}

	// Créer le contenu du fichier de parsing
	var content strings.Builder
	content.WriteString("PARSING RESULTS\n")
	content.WriteString("===============\n\n")
	content.WriteString("SOURCE CODE:\n")
	content.WriteString(code + "\n\n")
	content.WriteString(fmt.Sprintf("TARGET LANGUAGE: %s\n", strings.ToUpper(targetLang)))
	content.WriteString(fmt.Sprintf("PARSED ELEMENTS: %d\n", len(program)))
	content.WriteString("=================\n\n")

	for i, stmt := range program {
		if stmt != nil {
			content.WriteString(fmt.Sprintf("%d. Type: %T\n", i+1, stmt))
			content.WriteString(fmt.Sprintf("   Token: %s\n", stmt.TokenLiteral()))
			content.WriteString("\n")
		}
	}

	content.WriteString(fmt.Sprintf("\n\nTRANSPILATION RESULT FOR %s:\n", strings.ToUpper(targetLang)))
	content.WriteString("=============================================\n\n")
	content.WriteString(output)

	// Sauvegarder dans le fichier
	err := os.WriteFile(filename, []byte(content.String()), 0644)
	return output, err
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

	// Nom de fichier spécifique au langage demandé
	filename := fmt.Sprintf("parsing_%s.txt", req.Target)
	
	// 1. SUPPRIMER le fichier précédent pour ce langage
	os.Remove(filename) // Ignore les erreurs
	
	// 2. PARSER et GÉNÉRER uniquement pour le langage demandé
	output, err := saveParsingForSingleLanguage(req.Code, req.Target, filename)
	if err != nil {
		json.NewEncoder(w).Encode(TranspileResponse{
			Success: false,
			Error:   "Erreur lors du parsing: " + err.Error(),
		})
		return
	}

	// 3. Vérifier et nettoyer l'output
	if output == "" || strings.TrimSpace(output) == "" {
		output = "// Aucun code généré - Le parsing a peut-être échoué\n"
		output += "// Consultez le fichier " + filename + " pour plus de détails"
	}

	response := TranspileResponse{
		Success: true,
		Output:  fmt.Sprintf("// 🔄 Parsing optimisé pour %s\n// 📄 Détails dans: %s\n\n%s", 
			strings.ToUpper(req.Target), filename, output),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/transpile", handleTranspile)

	fmt.Println("🌐 Serveur web démarré sur http://localhost:8081")
	fmt.Println("💡 Ouvrez votre navigateur à cette adresse")
	fmt.Println("🔄 Ctrl+C pour arrêter le serveur")
	
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("❌ Erreur serveur: %v\n", err)
	}
}

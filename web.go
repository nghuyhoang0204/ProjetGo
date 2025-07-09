package main

import (
	"ProjetGo/generator"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type TranspilerResult struct {
	SourceCode   string
	JavaScript   string
	ErrorMessage string
	ParseTime    string
}

// Enhanced HTML template with modern features
const htmlTemplate = `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🚀 Transpilateur Multi-Langages v2.0</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        :root {
            --primary-color: #667eea;
            --secondary-color: #764ba2;
            --bg-color: #f8f9fa;
            --text-color: #333;
            --border-color: #e9ecef;
            --code-bg: #2d3748;
            --success-color: #28a745;
            --error-color: #dc3545;
            --warning-color: #ffc107;
        }
        
        [data-theme="dark"] {
            --bg-color: #1a1a1a;
            --text-color: #e0e0e0;
            --border-color: #404040;
            --code-bg: #2d2d2d;
        }
        
        body {
            font-family: 'Segoe UI', system-ui, sans-serif;
            background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
            min-height: 100vh;
            color: var(--text-color);
        }
        
        .container {
            max-width: 1600px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 30px;
            margin-bottom: 20px;
            text-align: center;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            background: linear-gradient(45deg, var(--primary-color), var(--secondary-color));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        
        .header p {
            color: #666;
            font-size: 1.1em;
        }
        
        .controls {
            display: flex;
            gap: 15px;
            align-items: center;
            justify-content: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .theme-toggle {
            background: var(--primary-color);
            color: white;
            border: none;
            padding: 10px 15px;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
        }
        
        .main-content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .input-section, .output-section {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .section-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 2px solid var(--border-color);
        }
        
        .section-title {
            font-size: 1.2em;
            font-weight: 600;
            color: var(--text-color);
        }
        
        .file-controls {
            display: flex;
            gap: 10px;
        }
        
        .btn {
            background: var(--primary-color);
            color: white;
            border: none;
            padding: 8px 15px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            transition: all 0.3s;
        }
        
        .btn:hover {
            transform: translateY(-1px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
        }
        
        .btn-secondary {
            background: var(--success-color);
        }
        
        .btn-warning {
            background: var(--warning-color);
            color: #333;
        }
        
        textarea {
            width: 100%;
            height: 400px;
            padding: 15px;
            border: 2px solid var(--border-color);
            border-radius: 8px;
            font-family: 'Fira Code', 'Courier New', monospace;
            font-size: 14px;
            resize: vertical;
            background: var(--bg-color);
            color: var(--text-color);
            transition: border-color 0.3s;
        }
        
        textarea:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 10px rgba(102, 126, 234, 0.3);
        }
        
        .language-selector {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
            flex-wrap: wrap;
        }
        
        .lang-btn {
            background: var(--bg-color);
            border: 2px solid var(--border-color);
            color: var(--text-color);
            padding: 8px 15px;
            border-radius: 20px;
            cursor: pointer;
            font-size: 12px;
            transition: all 0.3s;
        }
        
        .lang-btn.active {
            background: var(--primary-color);
            color: white;
            border-color: var(--primary-color);
        }
        
        .lang-btn:hover {
            transform: translateY(-1px);
        }
        
        .output-container {
            background: var(--code-bg);
            border-radius: 8px;
            padding: 15px;
            height: 400px;
            overflow-y: auto;
            font-family: 'Fira Code', 'Courier New', monospace;
            font-size: 13px;
            line-height: 1.4;
        }
        
        .output-container pre {
            margin: 0;
            color: #e2e8f0;
        }
        
        .status {
            padding: 10px 15px;
            border-radius: 8px;
            margin-top: 15px;
            font-weight: 500;
        }
        
        .status.success {
            background: rgba(40, 167, 69, 0.1);
            color: var(--success-color);
            border: 1px solid rgba(40, 167, 69, 0.3);
        }
        
        .status.error {
            background: rgba(220, 53, 69, 0.1);
            color: var(--error-color);
            border: 1px solid rgba(220, 53, 69, 0.3);
        }
        
        .status.info {
            background: rgba(102, 126, 234, 0.1);
            color: var(--primary-color);
            border: 1px solid rgba(102, 126, 234, 0.3);
        }
        
        .examples {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        
        .example-card {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 10px;
            padding: 15px;
            cursor: pointer;
            transition: all 0.3s;
            border: 2px solid transparent;
        }
        
        .example-card:hover {
            transform: translateY(-2px);
            border-color: var(--primary-color);
        }
        
        .example-title {
            font-weight: 600;
            margin-bottom: 8px;
            color: var(--text-color);
        }
        
        .example-desc {
            font-size: 12px;
            color: #666;
        }
        
        .loading {
            display: none;
            text-align: center;
            padding: 20px;
            color: var(--text-color);
        }
        
        .spinner {
            border: 3px solid var(--border-color);
            border-top: 3px solid var(--primary-color);
            border-radius: 50%;
            width: 30px;
            height: 30px;
            animation: spin 1s linear infinite;
            margin: 0 auto 10px;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
            }
            
            .controls {
                flex-direction: column;
            }
            
            .language-selector {
                justify-content: center;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Transpilateur Multi-Langages v2.0</h1>
            <p>Convertissez votre code TypeScript vers JavaScript</p>
        </div>
        
        <div class="controls">
            <button class="btn theme-toggle" onclick="toggleTheme()">🌙 Mode Sombre</button>
            <button class="btn btn-secondary" onclick="loadExample('basic')">📝 Exemple Basique</button>
            <button class="btn btn-secondary" onclick="loadExample('function')">🔧 Exemple Fonction</button>
            <button class="btn btn-secondary" onclick="loadExample('class')">🏗️ Exemple Classe</button>
            <button class="btn btn-warning" onclick="clearCode()">🗑️ Effacer</button>
        </div>
        
        <div class="main-content">
            <div class="input-section">
                <div class="section-header">
                    <div class="section-title">📝 Code Source (TypeScript)</div>
                    <div class="file-controls">
                        <input type="file" id="fileInput" accept=".ts,.js,.txt" style="display: none;" onchange="loadFile(event)">
                        <button class="btn" onclick="document.getElementById('fileInput').click()">📁 Ouvrir</button>
                        <button class="btn" onclick="downloadCode()">💾 Sauvegarder</button>
                    </div>
                </div>
                <textarea id="sourceCode" placeholder="Entrez votre code TypeScript ici...
Exemple :
const message: string = 'Hello World';
let count: number = 42;
const pi: number = 3.14;

function greet(name: string): string {
    return 'Hello ' + name;
}">{{.SourceCode}}</textarea>
            </div>
            
            <div class="output-section">
                <div class="section-header">
                    <div class="section-title">⚡ Code Généré</div>
                    <button class="btn" onclick="transpile()">🔄 Transpiler</button>
                </div>
                
                <div class="language-selector">
                    <button class="lang-btn active" data-lang="javascript">🟨 JavaScript</button>
                </div>
                
                <div class="loading" id="loading">
                    <div class="spinner"></div>
                    <div>Transpilation en cours...</div>
                </div>
                
                <div class="output-container" id="output">
                    <pre>Sélectionnez un langage cible et cliquez sur "Transpiler"</pre>
                </div>
                
                <div class="status info" id="status">Prêt à transpiler</div>
            </div>
        </div>
        
        <div class="examples">
            <div class="example-card" onclick="loadExample('basic')">
                <div class="example-title">🔤 Variables et Types</div>
                <div class="example-desc">Déclarations de variables avec types TypeScript</div>
            </div>
            <div class="example-card" onclick="loadExample('function')">
                <div class="example-title">🔧 Fonctions</div>
                <div class="example-desc">Fonctions avec paramètres typés et valeurs de retour</div>
            </div>
            <div class="example-card" onclick="loadExample('class')">
                <div class="example-title">🏗️ Classes et Interfaces</div>
                <div class="example-desc">Classes TypeScript avec méthodes et propriétés</div>
            </div>
            <div class="example-card" onclick="loadExample('advanced')">
                <div class="example-title">🚀 Fonctionnalités Avancées</div>
                <div class="example-desc">Template literals, arrays, objets et plus</div>
            </div>
        </div>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>
    
    <script>
        let currentTheme = 'light';
        let currentResults = {};
        
        // Theme toggle
        function toggleTheme() {
            currentTheme = currentTheme === 'light' ? 'dark' : 'light';
            document.body.setAttribute('data-theme', currentTheme === 'dark' ? 'dark' : 'light');
            document.querySelector('.theme-toggle').textContent = currentTheme === 'dark' ? '☀️ Mode Clair' : '🌙 Mode Sombre';
        }
        
        // Language selector
        document.querySelectorAll('.lang-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                document.querySelectorAll('.lang-btn').forEach(b => b.classList.remove('active'));
                this.classList.add('active');
                updateOutput();
            });
        });
        
        // Examples
                 const examples = {
             basic: 'const message: string = "Hello World";\nlet count: number = 42;\nconst pi: number = 3.14;\nlet isActive: boolean = true;',
            
                         function: 'function greet(name: string): string {\n    return "Hello " + name;\n}\n\nfunction add(a: number, b: number): number {\n    return a + b;\n}\n\nconst result = add(5, 3);\nconsole.log(greet("Alice"));',
             
             class: 'interface User {\n    id: number;\n    name: string;\n    email: string;\n}\n\nclass Calculator {\n    private value: number = 0;\n    \n    add(x: number): void {\n        this.value += x;\n    }\n    \n    getResult(): number {\n        return this.value;\n    }\n}\n\nconst calc = new Calculator();\ncalc.add(10);\nconsole.log(calc.getResult());',
             
             advanced: 'const users: User[] = [\n    { id: 1, name: "Alice", email: "alice@example.com" },\n    { id: 2, name: "Bob", email: "bob@example.com" }\n];\n\nconst template = "Hello " + users[0].name + "!";\nconst numbers: number[] = [1, 2, 3, 4, 5];\n\nfor (let i = 0; i < numbers.length; i++) {\n    console.log(numbers[i]);\n}'
        };
        
        function loadExample(type) {
            document.getElementById('sourceCode').value = examples[type] || examples.basic;
            transpile();
        }
        
        function clearCode() {
            document.getElementById('sourceCode').value = '';
            document.getElementById('output').innerHTML = '<pre>Sélectionnez un langage cible et cliquez sur "Transpiler"</pre>';
            document.getElementById('status').textContent = 'Code effacé';
            document.getElementById('status').className = 'status info';
        }
        
        // File operations
        function loadFile(event) {
            const file = event.target.files[0];
            if (file) {
                const reader = new FileReader();
                reader.onload = function(e) {
                    document.getElementById('sourceCode').value = e.target.result;
                    transpile();
                };
                reader.readAsText(file);
            }
        }
        
        function downloadCode() {
            const code = document.getElementById('sourceCode').value;
            if (code.trim()) {
                const blob = new Blob([code], { type: 'text/plain' });
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'code.ts';
                a.click();
                URL.revokeObjectURL(url);
            }
        }
        
        // Real-time transpilation
        let transpileTimeout;
        document.getElementById('sourceCode').addEventListener('input', function() {
            clearTimeout(transpileTimeout);
            transpileTimeout = setTimeout(transpile, 1000); // Debounce 1 second
        });
        
        async function transpile() {
            const code = document.getElementById('sourceCode').value;
            const activeLang = document.querySelector('.lang-btn.active').dataset.lang;
            
            if (!code.trim()) {
                document.getElementById('output').innerHTML = '<pre>Sélectionnez un langage cible et cliquez sur "Transpiler"</pre>';
                document.getElementById('status').textContent = 'Prêt à transpiler';
                document.getElementById('status').className = 'status info';
                return;
            }
            
            // Show loading
            document.getElementById('loading').style.display = 'block';
            document.getElementById('output').style.display = 'none';
            document.getElementById('status').textContent = 'Transpilation en cours...';
            document.getElementById('status').className = 'status info';
            
            try {
                const response = await fetch('/transpile', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        code: code,
                        target: activeLang
                    })
                });
                
                const result = await response.json();
                
                if (result.success) {
                    currentResults = result;
                    updateOutput();
                                         document.getElementById('status').textContent = 'Transpilation réussie (' + result.parseTime + ')';
                    document.getElementById('status').className = 'status success';
                } else {
                                         document.getElementById('output').innerHTML = '<pre class="error">' + result.error + '</pre>';
                    document.getElementById('status').textContent = 'Erreur de transpilation';
                    document.getElementById('status').className = 'status error';
                }
            } catch (error) {
                document.getElementById('output').innerHTML = '<pre class="error">Erreur de connexion</pre>';
                document.getElementById('status').textContent = 'Erreur de connexion';
                document.getElementById('status').className = 'status error';
            } finally {
                document.getElementById('loading').style.display = 'none';
                document.getElementById('output').style.display = 'block';
            }
        }
        
        function updateOutput() {
            const activeLang = document.querySelector('.lang-btn.active').dataset.lang;
            let output = '';
            
            if (activeLang === 'all') {
                const languages = [
                    { key: 'javascript', name: '🟨 JavaScript', code: currentResults.javascript },
                    { key: 'java', name: '☕ Java', code: currentResults.java },
                    { key: 'python', name: '🐍 Python', code: currentResults.python },
                    { key: 'csharp', name: '🔵 C#', code: currentResults.csharp },
                    { key: 'go', name: '🐹 Go', code: currentResults.go },
                    { key: 'rust', name: '🦀 Rust', code: currentResults.rust },
                    { key: 'swift', name: '🍎 Swift', code: currentResults.swift },
                    { key: 'php', name: '🐘 PHP', code: currentResults.php }
                ];
                
                languages.forEach(lang => {
                    if (lang.code) {
                                                 output += '<h4>' + lang.name + '</h4><pre><code class="language-' + lang.key + '">' + lang.code + '</code></pre>';
                    }
                });
            } else {
                const code = currentResults[activeLang];
                if (code) {
                                         output = '<pre><code class="language-' + activeLang + '">' + code + '</code></pre>';
                } else {
                    output = '<pre>Aucun code généré pour ce langage</pre>';
                }
            }
            
            document.getElementById('output').innerHTML = output;
            
            // Apply syntax highlighting
            if (window.Prism) {
                Prism.highlightAll();
            }
        }
        
        // Auto-transpile on page load if there's code
        window.addEventListener('load', function() {
            const code = document.getElementById('sourceCode').value;
            if (code.trim()) {
                transpile();
            }
        });
    </script>
</body>
</html>
`

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("transpiler").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := TranspilerResult{}

	if r.Method == "POST" {
		sourceCode := r.FormValue("source")
		result.SourceCode = sourceCode

		if sourceCode != "" {
			start := time.Now()

			// Transpilation
			l := lexer.New(sourceCode)
			p := parser.New(l)
			program := p.ParseProgram()

			elapsed := time.Since(start)
			result.ParseTime = elapsed.String()

			// Check for parsing errors
			if len(program.Statements) == 0 {
				result.ErrorMessage = "Aucun code valide détecté. Vérifiez la syntaxe de votre code source."
			} else {
				// Génération JavaScript uniquement
				result.JavaScript = generator.Generate(program.Statements)
			}
		}
	}

	tmpl.Execute(w, result)
}

// New API endpoint for real-time transpilation
func handleTranspile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Code   string `json:"code"`
		Target string `json:"target"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()

	l := lexer.New(request.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	elapsed := time.Since(start)

	response := map[string]interface{}{
		"success":   len(program.Statements) > 0,
		"parseTime": elapsed.String(),
	}

	if len(program.Statements) == 0 {
		response["error"] = "Aucun code valide détecté. Vérifiez la syntaxe."
	} else {
		// Utiliser la nouvelle fonction de transpilation directe
		response["javascript"] = generator.TranspileTS(request.Code)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StartWebServer() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/transpile", handleTranspile)

	fmt.Println("🌐 Serveur web démarré sur http://localhost:8080")
	fmt.Println("📝 Ouvrez votre navigateur et commencez à transpiler !")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

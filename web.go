package main

import (
	"html/template"
	"net/http"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
	"log"
)

type TranspilerResult struct {
	SourceCode   string
	JavaScript   string
	Java         string
	Python       string
	CSharp       string
	Go           string
	ErrorMessage string
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Transpilateur Multi-Langages</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
            font-size: 2.5em;
            background: linear-gradient(45deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .input-section {
            margin-bottom: 30px;
        }
        label {
            display: block;
            margin-bottom: 10px;
            font-weight: bold;
            color: #555;
            font-size: 1.1em;
        }
        textarea {
            width: 100%;
            height: 200px;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            resize: vertical;
            transition: border-color 0.3s;
        }
        textarea:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 10px rgba(102, 126, 234, 0.3);
        }
        .btn-container {
            text-align: center;
            margin: 20px 0;
        }
        button {
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            border: none;
            padding: 15px 30px;
            font-size: 16px;
            border-radius: 25px;
            cursor: pointer;
            transition: transform 0.3s, box-shadow 0.3s;
        }
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
        }
        .results {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
            gap: 20px;
            margin-top: 30px;
        }
        .result-box {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 10px;
            padding: 20px;
        }
        .result-box h3 {
            margin-top: 0;
            color: #495057;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .result-box pre {
            background: #2d3748;
            color: #e2e8f0;
            padding: 15px;
            border-radius: 8px;
            overflow-x: auto;
            font-family: 'Courier New', monospace;
            font-size: 13px;
            line-height: 1.4;
            margin: 0;
        }
        .error {
            background: #fee;
            border: 1px solid #fcc;
            color: #c33;
            padding: 15px;
            border-radius: 8px;
            margin: 10px 0;
        }
        .example-btn {
            background: #28a745;
            margin-left: 10px;
            padding: 10px 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Transpilateur Multi-Langages</h1>
        
        <form method="POST">
            <div class="input-section">
                <label for="source">Code Source (TypeScript-like) :</label>
                <textarea name="source" id="source" placeholder="Entrez votre code ici...
Exemple :
const message: string = &quot;Hello World&quot;;
let count: number = 42;
const pi: number = 3.14;">{{.SourceCode}}</textarea>
            </div>
            
            <div class="btn-container">
                <button type="submit">🔄 Transpiler</button>
                <button type="button" class="example-btn" onclick="loadExample()">📝 Exemple</button>
            </div>
        </form>

        {{if .ErrorMessage}}
        <div class="error">
            <strong>Erreur :</strong> {{.ErrorMessage}}
        </div>
        {{end}}

        {{if and (not .ErrorMessage) .JavaScript}}
        <div class="results">
            <div class="result-box">
                <h3>🟨 JavaScript</h3>
                <pre>{{.JavaScript}}</pre>
            </div>
            
            <div class="result-box">
                <h3>☕ Java</h3>
                <pre>{{.Java}}</pre>
            </div>
            
            <div class="result-box">
                <h3>🐍 Python</h3>
                <pre>{{.Python}}</pre>
            </div>
            
            <div class="result-box">
                <h3>🔷 C#</h3>
                <pre>{{.CSharp}}</pre>
            </div>
            
            <div class="result-box">
                <h3>🐹 Go</h3>
                <pre>{{.Go}}</pre>
            </div>
        </div>
        {{end}}
    </div>

    <script>
        function loadExample() {
            document.getElementById('source').value = 'const message: string = "Hello World";\nlet count: number = 42;\nconst pi: number = 3.14;\nlet isActive: boolean = true;';
        }
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
			// Transpilation
			l := lexer.New(sourceCode)
			p := parser.New(l)
			program := p.ParseProgram()

			// Génération dans tous les langages
			result.JavaScript = generator.Generate(program, generator.JavaScript)
			result.Java = generator.Generate(program, generator.Java)
			result.Python = generator.Generate(program, generator.Python)
			result.CSharp = generator.Generate(program, generator.CSharp)
			result.Go = generator.Generate(program, generator.Go)

			// Vérification si la transpilation a produit du contenu
			if result.JavaScript == "" && result.Java == "" {
				result.ErrorMessage = "Aucun code valide détecté. Vérifiez la syntaxe de votre code source."
			}
		}
	}

	tmpl.Execute(w, result)
}

func StartWebServer() {
	http.HandleFunc("/", handleHome)
	
	fmt.Println("🌐 Serveur web démarré sur http://localhost:8080")
	fmt.Println("📝 Ouvrez votre navigateur et commencez à transpiler !")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

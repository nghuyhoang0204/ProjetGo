package main

import (
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"ProjetGo/generator"
	"fmt"
	"strings"
)

func extractSimpleVariables(code string) string {
	// Pour l'instant, on va extraire manuellement les variables simples
	// qu'on peut trouver dans le code complexe
	
	lines := strings.Split(code, "\n")
	var extractedVars []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Chercher les déclarations de variables simples
		if strings.HasPrefix(line, "const ") || strings.HasPrefix(line, "let ") || strings.HasPrefix(line, "var ") {
			// Ignorer les lignes trop complexes (avec { ou [ ou fonction)
			if !strings.Contains(line, "{") && !strings.Contains(line, "[") && 
			   !strings.Contains(line, "(") && strings.Contains(line, "=") {
				extractedVars = append(extractedVars, line)
			}
		}
	}
	
	return strings.Join(extractedVars, "\n")
}

func mainSmartExtract() {
	complexCode := `type TaskStatus = 'pending' | 'in_progress' | 'done';

interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  createdAt: Date;
  updatedAt: Date;
}

class TaskManager {
  private tasks: Task[] = [];
  private nextId = 1;

  async createTask(title: string, description: string): Promise<Task> {
    const task: Task = {
      id: this.nextId++,
      title,
      description,
      status: 'pending',
      createdAt: new Date(),
      updatedAt: new Date(),
    };
    this.tasks.push(task);
    return task;
  }
}

const manager = new TaskManager();
const appName = "Task Manager";
let isRunning = true;
const version = "1.0.0";`

	fmt.Println("=== Code TypeScript Complexe ===")
	fmt.Println(complexCode)
	
	// Extraire les variables simples
	simpleVars := extractSimpleVariables(complexCode)
	fmt.Println("\n=== Variables Simples Extraites ===")
	fmt.Println(simpleVars)
	
	if simpleVars != "" {
		// Parser les variables simples
		l := lexer.New(simpleVars)
		p := parser.New(l)
		program := p.ParseProgram()
		
		fmt.Printf("\nVariables parsées: %d\n", len(program))
		
		// Générer en Java
		javaOutput := generator.Generate(program, generator.Java)
		fmt.Println("\n=== Java Output ===")
		fmt.Println(javaOutput)
		
		// Générer en Python
		pythonOutput := generator.Generate(program, generator.Python)
		fmt.Println("=== Python Output ===")
		fmt.Println(pythonOutput)
	} else {
		fmt.Println("Aucune variable simple trouvée à transpiler.")
	}
}

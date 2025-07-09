package main

import (
	"ProjetGo/ast"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"fmt"
	"strings"
)

// cleanTypeAnnotations nettoie les annotations de type pour avoir un AST JavaScript propre
func cleanTypeAnnotations(program *ast.Program) *ast.Program {
	cleanedStatements := make([]ast.Statement, 0, len(program.Statements))
	
	for _, stmt := range program.Statements {
		if stmt == nil {
			continue
		}
		
		// Filtrer les déclarations liées aux types
		switch s := stmt.(type) {
		case *ast.TypeAlias, *ast.Interface:
			// Ignorer complètement ces nœuds
			continue
		case *ast.ExpressionStatement:
			// Ignorer les expressions qui pourraient être des types
			if isTypeExpression(s.Expression) {
				continue
			}
			cleanedStatements = append(cleanedStatements, s)
		case *ast.FunctionDeclaration:
			// Garder la fonction mais nettoyer les annotations de type
			// Les types de paramètres et de retour sont déjà stockés séparément
			cleanedStatements = append(cleanedStatements, s)
		default:
			cleanedStatements = append(cleanedStatements, s)
		}
	}
	
	return &ast.Program{Statements: cleanedStatements}
}

// isTypeExpression détermine si une expression est probablement un type TypeScript
func isTypeExpression(expr ast.Expression) bool {
	if expr == nil {
		return false
	}

	// Les types comme "number", "string", etc. sont des identifiants
	if ident, ok := expr.(*ast.Identifier); ok {
		// Liste des types TypeScript communs
		typeNames := map[string]bool{
			"number": true,
			"string": true,
			"boolean": true,
			"any": true,
			"void": true,
			"never": true,
			"unknown": true,
			"object": true,
			"null": true,
			"undefined": true,
		}
		return typeNames[ident.Value]
	}
	
	return false
}

func TestCleanTranspilation() {
	// Code TypeScript à tester
	input := `
function addition(a: number, b: number): number {
  return a + b;
}
const resultat = addition(5, 3);
console.log("Résultat :", resultat); // Résultat : 8
`

	// Transpilation
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	
	// Nettoyer les annotations de type
	cleanedProgram := cleanTypeAnnotations(program)

	// Générer une sortie manuelle pour le débogage
	var sb strings.Builder
	
	// Fonction
	sb.WriteString("function addition(a, b) {\n")
	sb.WriteString("  return a + b;\n")
	sb.WriteString("}\n\n")
	
	// Utilisation
	sb.WriteString("const resultat = addition(5, 3);\n")
	sb.WriteString("console.log(\"Résultat :\", resultat);")
	
	expectedOutput := sb.String()

	// Afficher le résultat
	fmt.Println("TypeScript original:")
	fmt.Println("-------------------")
	fmt.Println(input)
	fmt.Println("\nJavaScript attendu:")
	fmt.Println("------------------")
	fmt.Println(expectedOutput)
}

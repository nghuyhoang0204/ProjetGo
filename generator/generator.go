package generator

import (
	"ProjetGo/ast"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"fmt"
	"strings"
)

// TypeScriptToJavaScriptTranspiler transpile TypeScript vers JavaScript
type TypeScriptToJavaScriptTranspiler struct{}

// Generate transpile un programme TypeScript vers JavaScript
func (t *TypeScriptToJavaScriptTranspiler) Generate(program *ast.Program) string {
	if program == nil {
		return ""
	}

	var result strings.Builder
	
	for _, statement := range program.Statements {
		jsCode := t.generateStatement(statement)
		if jsCode != "" {
			result.WriteString(jsCode)
			result.WriteString("\n")
		}
	}

	return result.String()
}

// generateStatement génère du JavaScript pour un statement
func (t *TypeScriptToJavaScriptTranspiler) generateStatement(stmt ast.Statement) string {
	switch node := stmt.(type) {
	case *ast.VariableDeclaration:
		return t.generateVariableDeclaration(node)
	case *ast.FunctionDeclaration:
		return t.generateFunctionDeclaration(node)
	case *ast.IfStatement:
		return t.generateIfStatement(node)
	case *ast.ForStatement:
		return t.generateForStatement(node)
	case *ast.WhileStatement:
		return t.generateWhileStatement(node)
	case *ast.BlockStatement:
		return t.generateBlockStatement(node)
	case *ast.ReturnStatement:
		return t.generateReturnStatement(node)
	case *ast.ExpressionStatement:
		return t.generateExpressionStatement(node)
	case *ast.AssignmentStatement:
		return t.generateAssignmentStatement(node)
	default:
		return ""
	}
}

// generateExpression génère du JavaScript pour une expression
func (t *TypeScriptToJavaScriptTranspiler) generateExpression(expr ast.Expression) string {
	switch node := expr.(type) {
	case *ast.Identifier:
		return node.Value
	case *ast.StringLiteral:
		return "\"" + node.Value + "\""
	case *ast.NumberLiteral:
		return node.Value
	case *ast.BooleanLiteral:
		if node.Value {
			return "true"
		}
		return "false"
	case *ast.ArrayLiteral:
		return t.generateArrayLiteral(node)
	case *ast.ObjectLiteral:
		return t.generateObjectLiteral(node)
	case *ast.CallExpression:
		return t.generateCallExpression(node)
	case *ast.MemberExpression:
		return t.generateMemberExpression(node)
	case *ast.InfixExpression:
		return t.generateInfixExpression(node)
	case *ast.TemplateLiteral:
		return t.generateTemplateLiteral(node)
	default:
		return ""
	}
}

// ============================================================================
// GENERATION DES STATEMENTS
// ============================================================================

func (t *TypeScriptToJavaScriptTranspiler) generateVariableDeclaration(node *ast.VariableDeclaration) string {
	var result strings.Builder

	// Utiliser const/let au lieu de var quand c'est approprié
	if node.IsConst {
		result.WriteString("const ")
	} else {
		result.WriteString("let ")
	}

	result.WriteString(node.Name)

	// Ignorer les types TypeScript (ils disparaissent en JavaScript)
	// if node.Type != "" { /* ignoré */ }

	if node.Value != nil {
		result.WriteString(" = ")
		result.WriteString(t.generateExpression(node.Value))
	}

	result.WriteString(";")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateFunctionDeclaration(node *ast.FunctionDeclaration) string {
	var result strings.Builder

	result.WriteString("function ")
	result.WriteString(node.Name)
	result.WriteString("(")

	// Paramètres (sans les types TypeScript)
	for i, param := range node.Parameters {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.Name)
		// Ignorer le type: param.Type
	}

	result.WriteString(")")
	// Ignorer le type de retour: node.ReturnType

	result.WriteString(" ")
	result.WriteString(t.generateBlockStatement(node.Body))

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateIfStatement(node *ast.IfStatement) string {
	var result strings.Builder

	result.WriteString("if (")
	result.WriteString(t.generateExpression(node.Condition))
	result.WriteString(") ")

	if blockStmt, ok := node.ThenBranch.(*ast.BlockStatement); ok {
		result.WriteString(t.generateBlockStatement(blockStmt))
	} else {
		result.WriteString(t.generateStatement(node.ThenBranch))
	}

	if node.ElseBranch != nil {
		result.WriteString(" else ")
		if blockStmt, ok := node.ElseBranch.(*ast.BlockStatement); ok {
			result.WriteString(t.generateBlockStatement(blockStmt))
		} else {
			result.WriteString(t.generateStatement(node.ElseBranch))
		}
	}

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateForStatement(node *ast.ForStatement) string {
	var result strings.Builder

	result.WriteString("for (")

	if node.Init != nil {
		initCode := t.generateStatement(node.Init)
		// Supprimer le point-virgule final pour l'init
		initCode = strings.TrimSuffix(initCode, ";")
		result.WriteString(initCode)
	}
	result.WriteString("; ")

	if node.Condition != nil {
		result.WriteString(t.generateExpression(node.Condition))
	}
	result.WriteString("; ")

	if node.Update != nil {
		updateCode := t.generateStatement(node.Update)
		// Supprimer le point-virgule final pour l'update
		updateCode = strings.TrimSuffix(updateCode, ";")
		result.WriteString(updateCode)
	}

	result.WriteString(") ")

	if blockStmt, ok := node.Body.(*ast.BlockStatement); ok {
		result.WriteString(t.generateBlockStatement(blockStmt))
	} else {
		result.WriteString(t.generateStatement(node.Body))
	}

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateWhileStatement(node *ast.WhileStatement) string {
	var result strings.Builder

	result.WriteString("while (")
	result.WriteString(t.generateExpression(node.Condition))
	result.WriteString(") ")

	if blockStmt, ok := node.Body.(*ast.BlockStatement); ok {
		result.WriteString(t.generateBlockStatement(blockStmt))
	} else {
		result.WriteString(t.generateStatement(node.Body))
	}

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateBlockStatement(node *ast.BlockStatement) string {
	var result strings.Builder

	result.WriteString("{\n")

	for _, statement := range node.Statements {
		jsCode := t.generateStatement(statement)
		if jsCode != "" {
			result.WriteString("    ") // Indentation
			result.WriteString(jsCode)
			result.WriteString("\n")
		}
	}

	result.WriteString("}")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateReturnStatement(node *ast.ReturnStatement) string {
	var result strings.Builder

	result.WriteString("return")

	if node.Value != nil {
		result.WriteString(" ")
		result.WriteString(t.generateExpression(node.Value))
	}

	result.WriteString(";")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateExpressionStatement(node *ast.ExpressionStatement) string {
	return t.generateExpression(node.Expression) + ";"
}

func (t *TypeScriptToJavaScriptTranspiler) generateAssignmentStatement(node *ast.AssignmentStatement) string {
	var result strings.Builder

	result.WriteString(node.Name)
	result.WriteString(" = ")
	result.WriteString(t.generateExpression(node.Value))
	result.WriteString(";")

	return result.String()
}

// ============================================================================
// GENERATION DES EXPRESSIONS
// ============================================================================

func (t *TypeScriptToJavaScriptTranspiler) generateArrayLiteral(node *ast.ArrayLiteral) string {
	var result strings.Builder

	result.WriteString("[")

	for i, element := range node.Elements {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(t.generateExpression(element))
	}

	result.WriteString("]")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateObjectLiteral(node *ast.ObjectLiteral) string {
	var result strings.Builder

	result.WriteString("{\n")

	for i, prop := range node.Properties {
		if i > 0 {
			result.WriteString(",\n")
		}
		result.WriteString("    ") // Indentation
		
		// Clé (avec ou sans guillemets selon le besoin)
		if t.needsQuotes(prop.Key) {
			result.WriteString("\"")
			result.WriteString(prop.Key)
			result.WriteString("\"")
		} else {
			result.WriteString(prop.Key)
		}
		
		result.WriteString(": ")
		result.WriteString(t.generateExpression(prop.Value))
	}

	result.WriteString("\n}")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateCallExpression(node *ast.CallExpression) string {
	var result strings.Builder

	result.WriteString(t.generateExpression(node.Function))
	result.WriteString("(")

	for i, arg := range node.Arguments {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(t.generateExpression(arg))
	}

	result.WriteString(")")
	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateMemberExpression(node *ast.MemberExpression) string {
	var result strings.Builder

	result.WriteString(t.generateExpression(node.Object))

	if node.Computed {
		result.WriteString("[")
		result.WriteString(t.generateExpression(node.Property))
		result.WriteString("]")
	} else {
		result.WriteString(".")
		result.WriteString(t.generateExpression(node.Property))
	}

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateInfixExpression(node *ast.InfixExpression) string {
	var result strings.Builder

	// Gérer les expressions préfixes (comme -x, !x)
	if node.Left == nil {
		result.WriteString(node.Operator)
		result.WriteString(t.generateExpression(node.Right))
		return result.String()
	}

	// Expression infixe normale
	result.WriteString(t.generateExpression(node.Left))
	result.WriteString(" ")
	result.WriteString(node.Operator)
	result.WriteString(" ")
	result.WriteString(t.generateExpression(node.Right))

	return result.String()
}

func (t *TypeScriptToJavaScriptTranspiler) generateTemplateLiteral(node *ast.TemplateLiteral) string {
	var result strings.Builder

	result.WriteString("`")
	
	for _, part := range node.Parts {
		if strLit, ok := part.(*ast.StringLiteral); ok {
			// Partie texte
			result.WriteString(strLit.Value)
		} else {
			// Expression interpolée
			result.WriteString("${")
			result.WriteString(t.generateExpression(part))
			result.WriteString("}")
		}
	}

	result.WriteString("`")
	return result.String()
}

// ============================================================================
// FONCTIONS UTILITAIRES
// ============================================================================

// needsQuotes détermine si une clé d'objet a besoin de guillemets
func (t *TypeScriptToJavaScriptTranspiler) needsQuotes(key string) bool {
	if key == "" {
		return true
	}

	// Vérifier si c'est un identifiant valide
	if !isValidIdentifier(key) {
		return true
	}

	// Vérifier si c'est un mot-clé réservé
	reservedWords := map[string]bool{
		"break": true, "case": true, "catch": true, "class": true, "const": true,
		"continue": true, "debugger": true, "default": true, "delete": true,
		"do": true, "else": true, "export": true, "extends": true, "finally": true,
		"for": true, "function": true, "if": true, "import": true, "in": true,
		"instanceof": true, "new": true, "return": true, "super": true, "switch": true,
		"this": true, "throw": true, "try": true, "typeof": true, "var": true,
		"void": true, "while": true, "with": true, "yield": true,
	}

	return reservedWords[key]
}

// isValidIdentifier vérifie si une chaîne est un identifiant valide
func isValidIdentifier(str string) bool {
	if str == "" {
		return false
	}

	// Premier caractère doit être une lettre, _, ou $
	first := str[0]
	if !(first >= 'a' && first <= 'z') && !(first >= 'A' && first <= 'Z') && first != '_' && first != '$' {
		return false
	}

	// Autres caractères peuvent être des lettres, chiffres, _ ou $
	for i := 1; i < len(str); i++ {
		ch := str[i]
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') && 
		   !(ch >= '0' && ch <= '9') && ch != '_' && ch != '$' {
			return false
		}
	}

	return true
}

// ============================================================================
// FONCTION PRINCIPALE DE TRANSPILATION
// ============================================================================

// TranspileTypeScriptToJavaScript transpile du code TypeScript vers JavaScript
func TranspileTypeScriptToJavaScript(typescriptCode string) (string, error) {
	// Créer le lexer
	l := lexer.New(typescriptCode)
	
	// Créer le parser
	p := parser.New(l)
	
	// Parser le programme
	program := p.ParseProgram()
	
	// Vérifier les erreurs de parsing
	if len(p.Errors()) > 0 {
		return "", fmt.Errorf("erreurs de parsing: %v", p.Errors())
	}
	
	// Générer le JavaScript
	transpiler := &TypeScriptToJavaScriptTranspiler{}
	javascriptCode := transpiler.Generate(program)
	
	return javascriptCode, nil
}

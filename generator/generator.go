package generator

import (
	"ProjetGo/ast"
	"ProjetGo/lexer"
	"ProjetGo/parser"
	"strings"
)

// TypeScriptToJavaScriptGenerator est le générateur de code pour JavaScript
type TypeScriptToJavaScriptGenerator struct{}

// Generate génère du code JavaScript à partir de l'AST TypeScript
func Generate(statements []ast.Statement) string {
	generator := &TypeScriptToJavaScriptGenerator{}
	generatedCode := generator.Generate(statements)
	
	// Nettoyer la sortie pour éliminer les artefacts typiques
	return CleanJavaScriptOutput(generatedCode)
}

// Generate génère du code JavaScript à partir d'un tableau de statements
func (g *TypeScriptToJavaScriptGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder
	var lastWasFunctionOrClass bool = false

	for _, stmt := range statements {
		// Ignorer les nœuds nil ou vides
		if stmt == nil {
			continue
		}
		
		// Générer le code en fonction du type de statement
		var code string
		
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			code = g.generateVariableDeclaration(s)
		case *ast.FunctionDeclaration:
			code = g.generateFunction(s)
			lastWasFunctionOrClass = true
		case *ast.IfStatement:
			code = g.generateIfStatement(s)
		case *ast.ForStatement:
			code = g.generateForStatement(s)
		case *ast.WhileStatement:
			code = g.generateWhileStatement(s)
		case *ast.ReturnStatement:
			code = g.generateReturnStatement(s)
		case *ast.ExpressionStatement:
			// Ignorer les expressions qui représentent des types (comme les identifiants "number", etc.)
			if isTypeExpression(s.Expression) {
				continue
			}
			code = g.generateExpressionStatement(s)
		case *ast.TypeAlias:
			// TypeScript: type aliases sont ignorés en JavaScript
			continue
		case *ast.Interface:
			// TypeScript: interfaces sont ignorées en JavaScript
			continue
		case *ast.ClassDeclaration:
			code = g.generateClass(s)
			lastWasFunctionOrClass = true
		}
		
		// N'ajouter que le code non vide
		if code != "" && code != ";" {
			// Ajouter une ligne vide après les fonctions et classes
			if lastWasFunctionOrClass {
				sb.WriteString("\n")
				lastWasFunctionOrClass = false
			}
			
			sb.WriteString(code)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// generateIfStatement génère le code pour une instruction if
func (g *TypeScriptToJavaScriptGenerator) generateIfStatement(is *ast.IfStatement) string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(g.generateExpression(is.Condition))
	sb.WriteString(") ")
	sb.WriteString(g.generateStatement(is.ThenBranch))
	if is.ElseBranch != nil {
		sb.WriteString(" else ")
		sb.WriteString(g.generateStatement(is.ElseBranch))
	}
	return sb.String()
}

// generateForStatement génère le code pour une boucle for
func (g *TypeScriptToJavaScriptGenerator) generateForStatement(fs *ast.ForStatement) string {
	var sb strings.Builder
	sb.WriteString("for (")
	
	// Initialisation
	if fs.Init != nil {
		sb.WriteString(g.generateStatement(fs.Init))
	} else {
		sb.WriteString("; ")
	}

	// Condition
	if fs.Condition != nil {
		sb.WriteString(g.generateExpression(fs.Condition))
	}
	sb.WriteString("; ")

	// Update
	if fs.Update != nil {
		// Ne pas ajouter de point-virgule car c'est une expression
		expression, ok := fs.Update.(*ast.ExpressionStatement)
		if ok {
			sb.WriteString(g.generateExpression(expression.Expression))
		} else {
			sb.WriteString(g.generateStatement(fs.Update))
		}
	}
	
	sb.WriteString(") ")
	sb.WriteString(g.generateStatement(fs.Body))
	
	return sb.String()
}

// generateWhileStatement génère le code pour une boucle while
func (g *TypeScriptToJavaScriptGenerator) generateWhileStatement(ws *ast.WhileStatement) string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(g.generateExpression(ws.Condition))
	sb.WriteString(") ")
	sb.WriteString(g.generateStatement(ws.Body))
	return sb.String()
}

// generateReturnStatement génère le code pour un return
func (g *TypeScriptToJavaScriptGenerator) generateReturnStatement(rs *ast.ReturnStatement) string {
	var sb strings.Builder
	sb.WriteString("return")
	
	if rs.Value != nil {
		sb.WriteString(" ")
		sb.WriteString(g.generateExpression(rs.Value))
	}
	
	sb.WriteString(";")
	return sb.String()
}

// generateExpressionStatement génère le code pour une expression utilisée comme instruction
func (g *TypeScriptToJavaScriptGenerator) generateExpressionStatement(es *ast.ExpressionStatement) string {
	// Ignorer complètement les expressions de type
	if isTypeExpression(es.Expression) {
		return ""
	}
	return g.generateExpression(es.Expression) + ";";
}

// generateStatement génère le code pour une instruction
func (g *TypeScriptToJavaScriptGenerator) generateStatement(stmt ast.Statement) string {
	if stmt == nil {
		return ""
	}

	switch s := stmt.(type) {
	case *ast.BlockStatement:
		return g.generateBlockStatement(s)
	case *ast.VariableDeclaration:
		return g.generateVariableDeclaration(s)
	case *ast.FunctionDeclaration:
		return g.generateFunction(s)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(s)
	case *ast.IfStatement:
		return g.generateIfStatement(s)
	case *ast.ForStatement:
		return g.generateForStatement(s)
	case *ast.WhileStatement:
		return g.generateWhileStatement(s)
	case *ast.ExpressionStatement:
		return g.generateExpressionStatement(s)
	default:
		return ""
	}
}

// generateBlockStatement génère le code pour un bloc d'instructions
func (g *TypeScriptToJavaScriptGenerator) generateBlockStatement(block *ast.BlockStatement) string {
	if block == nil || len(block.Statements) == 0 {
		return "{\n}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	
	for _, stmt := range block.Statements {
		// Ignorer les déclarations de type
		if stmt == nil {
			continue
		}
		
		// Ignorer les expressions qui sont probablement des types
		if es, ok := stmt.(*ast.ExpressionStatement); ok && isTypeExpression(es.Expression) {
			continue
		}
		
		// Générer le code pour cette instruction
		stmtCode := g.generateStatement(stmt)
		// Ne pas ajouter les instructions vides (comme celles qui ont été filtrées)
		if stmtCode != "" {
			sb.WriteString("  ")
			sb.WriteString(stmtCode)
			sb.WriteString("\n")
		}
	}
	
	sb.WriteString("}")
	return sb.String()
}

// generateFunction génère le code pour une fonction
func (g *TypeScriptToJavaScriptGenerator) generateFunction(fd *ast.FunctionDeclaration) string {
	var sb strings.Builder
	
	// Préfixe async
	if fd.IsAsync {
		sb.WriteString("async ")
	}
	
	// Fonction elle-même
	sb.WriteString("function ")
	sb.WriteString(fd.Name)
	sb.WriteString("(")
	
	// Paramètres - uniquement les noms, pas les types
	for i, param := range fd.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		
		// Ajouter les valeurs par défaut si présentes
		if param.DefaultValue != nil {
			sb.WriteString(" = ")
			sb.WriteString(g.generateExpression(param.DefaultValue))
		}
	}
	
	sb.WriteString(") ")
	// Le type de retour est complètement ignoré
	
	// Générer le corps de la fonction en filtrant les annotations de type
	if fd.Body != nil {
		// Filtrer le corps pour éliminer les expressions de type
		cleanedBody := &ast.BlockStatement{
			Token:      fd.Body.Token,
			Statements: []ast.Statement{},
		}
		
		for _, stmt := range fd.Body.Statements {
			if stmt == nil {
				continue
			}
			
			if es, ok := stmt.(*ast.ExpressionStatement); ok {
				if isTypeExpression(es.Expression) {
					continue
				}
			}
			
			cleanedBody.Statements = append(cleanedBody.Statements, stmt)
		}
		
		sb.WriteString(g.generateBlockStatement(cleanedBody))
	} else {
		sb.WriteString("{\n}")
	}
	
	return sb.String()
}

// generateVariableDeclaration génère le code pour une déclaration de variable
func (g *TypeScriptToJavaScriptGenerator) generateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	// En JavaScript on conserve let/const mais on supprime var
	if vd.IsConst {
		sb.WriteString("const ")
	} else {
		sb.WriteString("let ")
	}
	
	sb.WriteString(vd.Name)
	
	// Ignorer les types TypeScript (ils disparaissent en JavaScript)
	
	if vd.Value != nil {
		sb.WriteString(" = ")
		sb.WriteString(g.generateExpression(vd.Value))
	}
	
	sb.WriteString(";")
	return sb.String()
}

// generateClass génère le code pour une classe
func (g *TypeScriptToJavaScriptGenerator) generateClass(cd *ast.ClassDeclaration) string {
	var sb strings.Builder
	sb.WriteString("class ")
	sb.WriteString(cd.Name)
	
	if cd.SuperClass != "" {
		sb.WriteString(" extends ")
		sb.WriteString(cd.SuperClass)
	}
	
	sb.WriteString(" {\n")
	
	// Propriétés et méthodes
	for _, prop := range cd.Properties {
		sb.WriteString("  ")
		// Ignorer les types en JavaScript et les modificateurs private/readonly
		if prop.Value != nil {
			sb.WriteString(prop.Name)
			sb.WriteString(" = ")
			sb.WriteString(g.generateExpression(prop.Value))
			sb.WriteString(";\n")
		}
	}
	
	for _, method := range cd.Methods {
		sb.WriteString("  ")
		if method.IsStatic {
			sb.WriteString("static ")
		}
		// Ignorer private en JavaScript
		
		sb.WriteString(method.Name)
		sb.WriteString("(")
		
		// Paramètres
		for i, param := range method.Parameters {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(param.Name)
			
			// Ajouter la valeur par défaut si elle existe
			if param.DefaultValue != nil {
				sb.WriteString(" = ")
				sb.WriteString(g.generateExpression(param.DefaultValue))
			}
			// Ignorer les types en JavaScript
		}
		
		sb.WriteString(") ")
		// Ignorer le type de retour
		sb.WriteString(g.generateBlockStatement(method.Body))
		sb.WriteString("\n")
	}
	
	sb.WriteString("}")
	return sb.String()
}

// generateExpression génère le code pour une expression
func (g *TypeScriptToJavaScriptGenerator) generateExpression(expr ast.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.StringLiteral:
		// S'assurer que la chaîne a bien ses délimiteurs
		if !strings.HasPrefix(e.Value, "\"") && !strings.HasPrefix(e.Value, "'") {
			return "\"" + e.Value + "\""
		}
		return e.Value
	case *ast.NumberLiteral:
		return e.Value
	case *ast.BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.ArrayLiteral:
		return g.generateArrayLiteral(e)
	case *ast.ObjectLiteral:
		return g.generateObjectLiteral(e)
	case *ast.TemplateLiteral:
		return g.generateTemplateLiteral(e)
	case *ast.CallExpression:
		return g.generateCallExpression(e)
	case *ast.MemberExpression:
		return g.generateMemberExpression(e)
	case *ast.InfixExpression:
		return g.generateInfixExpression(e)
	default:
		return ""
	}
}

// generateArrayLiteral génère le code pour un tableau
func (g *TypeScriptToJavaScriptGenerator) generateArrayLiteral(al *ast.ArrayLiteral) string {
	var sb strings.Builder
	sb.WriteString("[")
	
	for i, element := range al.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.generateExpression(element))
	}
	
	sb.WriteString("]")
	return sb.String()
}

// generateObjectLiteral génère le code pour un objet
func (g *TypeScriptToJavaScriptGenerator) generateObjectLiteral(ol *ast.ObjectLiteral) string {
	var sb strings.Builder
	sb.WriteString("{")
	
	for i, property := range ol.Properties {
		if i > 0 {
			sb.WriteString(", ")
		}
		
		sb.WriteString(property.Key)
		sb.WriteString(": ")
		sb.WriteString(g.generateExpression(property.Value))
	}
	
	sb.WriteString("}")
	return sb.String()
}

// generateTemplateLiteral génère le code pour un template string
func (g *TypeScriptToJavaScriptGenerator) generateTemplateLiteral(tl *ast.TemplateLiteral) string {
	var sb strings.Builder
	
	for _, part := range tl.Parts {
		sb.WriteString(g.generateExpression(part))
	}
	
	return sb.String()
}

// generateCallExpression génère le code pour un appel de fonction
func (g *TypeScriptToJavaScriptGenerator) generateCallExpression(ce *ast.CallExpression) string {
	var sb strings.Builder
	sb.WriteString(g.generateExpression(ce.Function))
	sb.WriteString("(")
	
	for i, arg := range ce.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.generateExpression(arg))
	}
	
	sb.WriteString(")")
	return sb.String()
}

// generateMemberExpression génère le code pour un accès à une propriété
func (g *TypeScriptToJavaScriptGenerator) generateMemberExpression(me *ast.MemberExpression) string {
	var sb strings.Builder
	sb.WriteString(g.generateExpression(me.Object))
	
	if me.Computed {
		sb.WriteString("[")
		sb.WriteString(g.generateExpression(me.Property))
		sb.WriteString("]")
	} else {
		sb.WriteString(".")
		sb.WriteString(g.generateExpression(me.Property))
	}
	
	return sb.String()
}

// generateInfixExpression génère le code pour une expression avec opérateur
func (g *TypeScriptToJavaScriptGenerator) generateInfixExpression(ie *ast.InfixExpression) string {
	var sb strings.Builder
	
	// Gestion des cas particuliers (opérateurs unaires, etc.)
	if ie.Left == nil {
		// C'est probablement un opérateur préfixe (unaire)
		sb.WriteString(ie.Operator)
		sb.WriteString(g.generateExpression(ie.Right))
	} else {
		// Opérateur infixe normal
		sb.WriteString(g.generateExpression(ie.Left))
		sb.WriteString(" ")
		sb.WriteString(ie.Operator)
		sb.WriteString(" ")
		sb.WriteString(g.generateExpression(ie.Right))
	}
	
	return sb.String()
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
			"Array": true,
			"Promise": true,
			"Record": true,
			"Partial": true,
			// Autres génériques courants
			"T": true,
			"U": true,
			"V": true,
			"K": true,
		}
		return typeNames[ident.Value]
	}
	
	// Les expressions infixes qui contiennent des opérateurs de type (comme | ou &) sont aussi des types
	if infix, ok := expr.(*ast.InfixExpression); ok {
		if infix.Operator == "|" || infix.Operator == "&" || infix.Operator == ":" {
			return true
		}
		// Si l'un des côtés est un type, alors c'est probablement une expression de type
		return isTypeExpression(infix.Left) || isTypeExpression(infix.Right)
	}

	// Les appels d'expressions sont des types génériques si leur fonction est un type
	if call, ok := expr.(*ast.CallExpression); ok {
		return isTypeExpression(call.Function)
	}
	
	return false
}

// TranspileTS transpile du code TypeScript en JavaScript
// Cette fonction peut être utilisée directement avec le code source TypeScript
func TranspileTS(typescriptCode string) string {
	// Méthode 1 : Utiliser le parser/lexer existant puis nettoyer
	l := lexer.New(typescriptCode)
	p := parser.New(l)
	program := p.ParseProgram()
	
	if program == nil || len(program.Statements) == 0 {
		// Si le parsing échoue, utiliser la méthode directe
		return GenerateFromSource(typescriptCode)
	}
	
	// Générer le code JavaScript avec le générateur standard
	generatedCode := Generate(program.Statements)
	
	// Si le résultat est trop petit ou semble cassé, utiliser la méthode directe
	if len(generatedCode) < 10 || !strings.Contains(generatedCode, "{") {
		return GenerateFromSource(typescriptCode)
	}
	
	// Sinon, retourner la version nettoyée du générateur standard
	return generatedCode
}

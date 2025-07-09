package generator

import (
	"strings"
	"ProjetGo/ast"
)

// TargetLanguage représente les langages de sortie supportés
type TargetLanguage string

const (
	JavaScript TargetLanguage = "javascript"
	Java       TargetLanguage = "java"
	Python     TargetLanguage = "python"
	CSharp     TargetLanguage = "csharp"
	Go         TargetLanguage = "go"
)

// CodeGenerator interface pour tous les générateurs de code
type CodeGenerator interface {
	Generate(statements []ast.Statement) string
}

// Generate génère du code dans le langage cible spécifié
func Generate(statements []ast.Statement, targetLang TargetLanguage) string {
	var generator CodeGenerator

	switch targetLang {
	case JavaScript:
		generator = &JavaScriptGenerator{}
	case Java:
		generator = &JavaGenerator{}
	case Python:
		generator = &PythonGenerator{}
	case CSharp:
		generator = &CSharpGenerator{}
	case Go:
		generator = &GoGenerator{}
	default:
		generator = &JavaScriptGenerator{} // défaut
	}

	return generator.Generate(statements)
}

// JavaScriptGenerator génère du code JavaScript
type JavaScriptGenerator struct{}

func (jsg *JavaScriptGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			sb.WriteString(jsg.GenerateVariableDeclaration(s))
		case *ast.FunctionDeclaration:
			sb.WriteString(jsg.GenerateFunction(s))
		case *ast.IfStatement:
			sb.WriteString(jsg.GenerateIfStatement(s))
		case *ast.ForStatement:
			sb.WriteString(jsg.GenerateForStatement(s))
		case *ast.WhileStatement:
			sb.WriteString(jsg.GenerateWhileStatement(s))
		case *ast.ReturnStatement:
			sb.WriteString(jsg.GenerateReturnStatement(s))
		case *ast.ExpressionStatement:
			sb.WriteString(jsg.GenerateExpressionStatement(s))
		case *ast.TypeAlias:
			sb.WriteString(jsg.GenerateTypeAlias(s))
		case *ast.Interface:
			sb.WriteString(jsg.GenerateInterface(s))
		case *ast.ClassDeclaration:
			sb.WriteString(jsg.GenerateClass(s))
		}
	}

	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateIfStatement(is *ast.IfStatement) string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(jsg.GenerateExpression(is.Condition))
	sb.WriteString(") ")
	sb.WriteString(jsg.GenerateStatement(is.ThenBranch))
	if is.ElseBranch != nil {
		sb.WriteString(" else ")
		sb.WriteString(jsg.GenerateStatement(is.ElseBranch))
	}
	sb.WriteString("\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateForStatement(fs *ast.ForStatement) string {
	var sb strings.Builder
	sb.WriteString("for (")
	if fs.Init != nil {
		sb.WriteString(strings.TrimSuffix(jsg.GenerateStatement(fs.Init), "\n"))
	}
	sb.WriteString("; ")
	if fs.Condition != nil {
		sb.WriteString(jsg.GenerateExpression(fs.Condition))
	}
	sb.WriteString("; ")
	if fs.Update != nil {
		sb.WriteString(strings.TrimSuffix(jsg.GenerateStatement(fs.Update), "\n"))
	}
	sb.WriteString(") ")
	sb.WriteString(jsg.GenerateStatement(fs.Body))
	sb.WriteString("\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateWhileStatement(ws *ast.WhileStatement) string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(jsg.GenerateExpression(ws.Condition))
	sb.WriteString(") ")
	sb.WriteString(jsg.GenerateStatement(ws.Body))
	sb.WriteString("\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateReturnStatement(rs *ast.ReturnStatement) string {
	var sb strings.Builder
	sb.WriteString("return")
	if rs.Value != nil {
		sb.WriteString(" ")
		sb.WriteString(jsg.GenerateExpression(rs.Value))
	}
	sb.WriteString(";\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateExpressionStatement(es *ast.ExpressionStatement) string {
	return jsg.GenerateExpression(es.Expression) + ";\n"
}

func (jsg *JavaScriptGenerator) GenerateStatement(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.BlockStatement:
		return jsg.GenerateBlockStatement(s)
	case *ast.VariableDeclaration:
		return jsg.GenerateVariableDeclaration(s)
	case *ast.ExpressionStatement:
		return jsg.GenerateExpressionStatement(s)
	case *ast.ReturnStatement:
		return jsg.GenerateReturnStatement(s)
	}
	return ""
}

func (jsg *JavaScriptGenerator) GenerateBlockStatement(bs *ast.BlockStatement) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, stmt := range bs.Statements {
		sb.WriteString("    ")
		sb.WriteString(jsg.GenerateStatement(stmt))
	}
	sb.WriteString("}")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return jsg.GenerateStringLiteral(e)
	case *ast.NumberLiteral:
		return jsg.GenerateNumberLiteral(e)
	case *ast.BooleanLiteral:
		return jsg.GenerateBooleanLiteral(e)
	case *ast.TemplateLiteral:
		return jsg.GenerateTemplateLiteral(e)
	case *ast.Identifier:
		return e.Value
	case *ast.InfixExpression:
		return jsg.GenerateExpression(e.Left) + " " + e.Operator + " " + jsg.GenerateExpression(e.Right)
	case *ast.ArrayLiteral:
		return jsg.GenerateArrayLiteral(e)
	case *ast.ObjectLiteral:
		return jsg.GenerateObjectLiteral(e)
	case *ast.CallExpression:
		return jsg.GenerateCallExpression(e)
	case *ast.IndexExpression:
		return jsg.GenerateIndexExpression(e)
	case *ast.DotExpression:
		return jsg.GenerateDotExpression(e)
	}
	return ""
}

func (jsg *JavaScriptGenerator) GenerateTemplateLiteral(tl *ast.TemplateLiteral) string {
	if len(tl.Parts) > 0 {
		return "`" + tl.Parts[0].TokenLiteral() + "`"
	}
	return "`template`"
}

func (jsg *JavaScriptGenerator) GenerateArrayLiteral(al *ast.ArrayLiteral) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, element := range al.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(jsg.GenerateExpression(element))
	}
	sb.WriteString("]")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateObjectLiteral(ol *ast.ObjectLiteral) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, prop := range ol.Properties {
		if i > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString("  ")
		sb.WriteString(prop.Key)
		sb.WriteString(": ")
		sb.WriteString(jsg.GenerateExpression(prop.Value))
	}
	sb.WriteString("\n}")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateCallExpression(ce *ast.CallExpression) string {
	var sb strings.Builder
	sb.WriteString(jsg.GenerateExpression(ce.Function))
	sb.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(jsg.GenerateExpression(arg))
	}
	sb.WriteString(")")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateIndexExpression(ie *ast.IndexExpression) string {
	return jsg.GenerateExpression(ie.Left) + "[" + jsg.GenerateExpression(ie.Index) + "]"
}

func (jsg *JavaScriptGenerator) GenerateDotExpression(de *ast.DotExpression) string {
	return jsg.GenerateExpression(de.Object) + "." + de.Property
}

func (jsg *JavaScriptGenerator) GenerateTypeAlias(ta *ast.TypeAlias) string {
	return "// Type alias: " + ta.Name + "\n"
}

func (jsg *JavaScriptGenerator) GenerateInterface(i *ast.Interface) string {
	return "// Interface: " + i.Name + "\n"
}

func (jsg *JavaScriptGenerator) GenerateClass(cd *ast.ClassDeclaration) string {
	var sb strings.Builder
	sb.WriteString("class ")
	sb.WriteString(cd.Name)
	sb.WriteString(" {\n")
	sb.WriteString("    // TODO: Implement class body\n")
	sb.WriteString("}\n\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateFunction(fd *ast.FunctionDeclaration) string {
	var sb strings.Builder
	
	if fd.IsAsync {
		sb.WriteString("async ")
	}
	sb.WriteString("function ")
	sb.WriteString(fd.Name)
	sb.WriteString("(")
	
	// Paramètres
	for i, param := range fd.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
	}
	
	sb.WriteString(") {\n")
	
	// Corps de la fonction
	for _, stmt := range fd.Body {
		sb.WriteString("    ")
		sb.WriteString(jsg.GenerateStatement(stmt))
	}
	
	sb.WriteString("}\n\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	if vd.IsConst {
		sb.WriteString("const ")
	} else {
		sb.WriteString("let ")
	}
	sb.WriteString(vd.Name)
	sb.WriteString(" = ")
	
	if vd.Value != nil {
		switch val := vd.Value.(type) {
		case *ast.StringLiteral:
			sb.WriteString(jsg.GenerateStringLiteral(val))
		case *ast.NumberLiteral:
			sb.WriteString(jsg.GenerateNumberLiteral(val))
		case *ast.BooleanLiteral:
			sb.WriteString(jsg.GenerateBooleanLiteral(val))
		case *ast.ArrayLiteral:
			sb.WriteString(jsg.GenerateArrayLiteral(val))
		case *ast.ObjectLiteral:
			sb.WriteString(jsg.GenerateObjectLiteral(val))
		case *ast.TemplateLiteral:
			sb.WriteString(jsg.GenerateTemplateLiteral(val))
		default:
			sb.WriteString(jsg.GenerateExpression(val))
		}
	}
	
	sb.WriteString(";\n")
	return sb.String()
}

func (jsg *JavaScriptGenerator) GenerateStringLiteral(sl *ast.StringLiteral) string {
	return "\"" + sl.Value + "\""
}

func (jsg *JavaScriptGenerator) GenerateNumberLiteral(nl *ast.NumberLiteral) string {
	return nl.Value
}

func (jsg *JavaScriptGenerator) GenerateBooleanLiteral(bl *ast.BooleanLiteral) string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// JavaGenerator génère du code Java
type JavaGenerator struct{}

func (jg *JavaGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder
	
	sb.WriteString("public class GeneratedCode {\n")
	
	// Séparer les variables et les fonctions
	var variables []ast.Statement
	var functions []ast.Statement
	var expressions []ast.Statement
	
	for _, stmt := range statements {
		switch stmt.(type) {
		case *ast.VariableDeclaration:
			variables = append(variables, stmt)
		case *ast.FunctionDeclaration:
			functions = append(functions, stmt)
		default:
			expressions = append(expressions, stmt)
		}
	}
	
	// Générer les fonctions d'abord
	for _, stmt := range functions {
		if s, ok := stmt.(*ast.FunctionDeclaration); ok {
			sb.WriteString("    ")
			sb.WriteString(jg.GenerateJavaFunction(s))
		}
	}
	
	// Méthode main
	sb.WriteString("    public static void main(String[] args) {\n")
	
	// Variables dans main
	for _, stmt := range variables {
		if s, ok := stmt.(*ast.VariableDeclaration); ok {
			sb.WriteString("        ")
			sb.WriteString(jg.GenerateVariableDeclaration(s))
		}
	}
	
	// Expressions/appels dans main  
	for _, stmt := range expressions {
		if s, ok := stmt.(*ast.ExpressionStatement); ok {
			sb.WriteString("        ")
			sb.WriteString(jg.GenerateJavaExpressionStatement(s))
		}
	}
	
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String()
}

func (jg *JavaGenerator) GenerateJavaFunction(fd *ast.FunctionDeclaration) string {
	var sb strings.Builder
	
	sb.WriteString("public static ")
	
	// Type de retour
	if fd.ReturnType == "void" || fd.ReturnType == "" {
		sb.WriteString("void ")
	} else if fd.ReturnType == "string" {
		sb.WriteString("String ")
	} else if fd.ReturnType == "number" {
		sb.WriteString("int ")
	} else if fd.ReturnType == "boolean" {
		sb.WriteString("boolean ")
	} else {
		sb.WriteString("Object ")
	}
	
	sb.WriteString(fd.Name)
	sb.WriteString("(")
	
	// Paramètres
	for i, param := range fd.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		
		// Type du paramètre
		if param.Type == "string" {
			sb.WriteString("String ")
		} else if param.Type == "number" {
			sb.WriteString("int ")
		} else if param.Type == "boolean" {
			sb.WriteString("boolean ")
		} else {
			sb.WriteString("Object ")
		}
		
		sb.WriteString(param.Name)
	}
	
	sb.WriteString(") {\n")
	
	// Corps de la fonction
	for _, stmt := range fd.Body {
		sb.WriteString("        ")
		sb.WriteString(jg.GenerateJavaStatement(stmt))
	}
	
	sb.WriteString("    }\n\n")
	return sb.String()
}

func (jg *JavaGenerator) GenerateJavaStatement(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		if s.Value != nil {
			return "return " + jg.GenerateExpression(s.Value) + ";\n"
		}
		return "return;\n"
	case *ast.IfStatement:
		return jg.GenerateJavaIfStatement(s)
	case *ast.VariableDeclaration:
		return jg.GenerateVariableDeclaration(s)
	}
	return ""
}

func (jg *JavaGenerator) GenerateJavaIfStatement(is *ast.IfStatement) string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(jg.GenerateExpression(is.Condition))
	sb.WriteString(") {\n")
	
	if blockStmt, ok := is.ThenBranch.(*ast.BlockStatement); ok {
		for _, stmt := range blockStmt.Statements {
			sb.WriteString("            ")
			sb.WriteString(jg.GenerateJavaStatement(stmt))
		}
	}
	
	sb.WriteString("        }")
	
	if is.ElseBranch != nil {
		sb.WriteString(" else {\n")
		if blockStmt, ok := is.ElseBranch.(*ast.BlockStatement); ok {
			for _, stmt := range blockStmt.Statements {
				sb.WriteString("            ")
				sb.WriteString(jg.GenerateJavaStatement(stmt))
			}
		}
		sb.WriteString("        }")
	}
	
	sb.WriteString("\n")
	return sb.String()
}

func (jg *JavaGenerator) GenerateJavaExpressionStatement(es *ast.ExpressionStatement) string {
	// Convertir console.log en System.out.println
	if callExpr, ok := es.Expression.(*ast.CallExpression); ok {
		if dotExpr, ok := callExpr.Function.(*ast.DotExpression); ok {
			if ident, ok := dotExpr.Object.(*ast.Identifier); ok {
				if ident.Value == "console" && dotExpr.Property == "log" {
					var sb strings.Builder
					sb.WriteString("System.out.println(")
					for i, arg := range callExpr.Arguments {
						if i > 0 {
							sb.WriteString(" + \", \" + ")
						}
						sb.WriteString(jg.GenerateExpression(arg))
					}
					sb.WriteString(");\n")
					return sb.String()
				}
			}
		}
	}
	
	return jg.GenerateExpression(es.Expression) + ";\n"
}

func (jg *JavaGenerator) GenerateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	// En Java, tout est final ou pas, pas de distinction const/let comme JS
	if vd.IsConst {
		sb.WriteString("final ")
	}
	
	// Déterminer le type Java
	switch vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString("String ")
	case *ast.NumberLiteral:
		sb.WriteString("int ") // Simplifié, pourrait être double/long
	case *ast.BooleanLiteral:
		sb.WriteString("boolean ")
	case *ast.ArrayLiteral:
		sb.WriteString("int[] ") // Simplifié pour les arrays de nombres
	case *ast.ObjectLiteral:
		sb.WriteString("java.util.HashMap<String, Object> ")
	case *ast.TemplateLiteral:
		sb.WriteString("String ")
	default:
		sb.WriteString("Object ")
	}
	
	sb.WriteString(vd.Name)
	sb.WriteString(" = ")
	
	if vd.Value != nil {
		switch val := vd.Value.(type) {
		case *ast.StringLiteral:
			sb.WriteString(jg.GenerateStringLiteral(val))
		case *ast.NumberLiteral:
			sb.WriteString(jg.GenerateNumberLiteral(val))
		case *ast.BooleanLiteral:
			sb.WriteString(jg.GenerateBooleanLiteral(val))
		case *ast.ArrayLiteral:
			sb.WriteString(jg.GenerateArrayLiteral(val))
		case *ast.ObjectLiteral:
			sb.WriteString(jg.GenerateObjectLiteral(val))
		case *ast.TemplateLiteral:
			sb.WriteString(jg.GenerateTemplateLiteral(val))
		default:
			sb.WriteString(jg.GenerateExpression(val))
		}
	}
	
	sb.WriteString(";\n")
	return sb.String()
}

func (jg *JavaGenerator) GenerateCallExpression(ce *ast.CallExpression) string {
	var sb strings.Builder
	sb.WriteString(jg.GenerateExpression(ce.Function))
	sb.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(jg.GenerateExpression(arg))
	}
	sb.WriteString(")")
	return sb.String()
}

func (jg *JavaGenerator) GenerateIndexExpression(ie *ast.IndexExpression) string {
	return jg.GenerateExpression(ie.Left) + "[" + jg.GenerateExpression(ie.Index) + "]"
}

func (jg *JavaGenerator) GenerateArrayLiteral(al *ast.ArrayLiteral) string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, element := range al.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(jg.GenerateExpression(element))
	}
	sb.WriteString("}")
	return sb.String()
}

func (jg *JavaGenerator) GenerateObjectLiteral(ol *ast.ObjectLiteral) string {
	// En Java, on va créer un HashMap ou une classe anonyme
	var sb strings.Builder
	sb.WriteString("new java.util.HashMap<String, Object>() {{")
	for i, prop := range ol.Properties {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(" put(\"")
		sb.WriteString(prop.Key)
		sb.WriteString("\", ")
		sb.WriteString(jg.GenerateExpression(prop.Value))
		sb.WriteString(")")
	}
	sb.WriteString("; }}")
	return sb.String()
}

func (jg *JavaGenerator) GenerateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return jg.GenerateStringLiteral(e)
	case *ast.NumberLiteral:
		return jg.GenerateNumberLiteral(e)
	case *ast.BooleanLiteral:
		return jg.GenerateBooleanLiteral(e)
	case *ast.TemplateLiteral:
		return jg.GenerateTemplateLiteral(e)
	case *ast.ArrayLiteral:
		return jg.GenerateArrayLiteral(e)
	case *ast.ObjectLiteral:
		return jg.GenerateObjectLiteral(e)
	case *ast.CallExpression:
		return jg.GenerateCallExpression(e)
	case *ast.IndexExpression:
		return jg.GenerateIndexExpression(e)
	case *ast.Identifier:
		return e.Value
	case *ast.InfixExpression:
		return jg.GenerateExpression(e.Left) + " " + e.Operator + " " + jg.GenerateExpression(e.Right)
	}
	return ""
}

func (jg *JavaGenerator) GenerateTemplateLiteral(tl *ast.TemplateLiteral) string {
	// En Java, convertir les template literals en String.format ou concaténation simple
	if len(tl.Parts) > 0 {
		content := tl.Parts[0].TokenLiteral()
		// Pour l'instant, retourner comme string simple (plus tard on peut parser ${} pour interpolation)
		return "\"" + content + "\""
	}
	return "\"\""
}

func (jg *JavaGenerator) GenerateStringLiteral(sl *ast.StringLiteral) string {
	return "\"" + sl.Value + "\""
}

func (jg *JavaGenerator) GenerateNumberLiteral(nl *ast.NumberLiteral) string {
	return nl.Value
}

func (jg *JavaGenerator) GenerateBooleanLiteral(bl *ast.BooleanLiteral) string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// PythonGenerator génère du code Python
type PythonGenerator struct{}

func (pg *PythonGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder

	// Séparer les variables et les fonctions
	var variables []ast.Statement
	var functions []ast.Statement
	var expressions []ast.Statement
	
	for _, stmt := range statements {
		switch stmt.(type) {
		case *ast.VariableDeclaration:
			variables = append(variables, stmt)
		case *ast.FunctionDeclaration:
			functions = append(functions, stmt)
		default:
			expressions = append(expressions, stmt)
		}
	}
	
	// Générer les fonctions d'abord
	for _, stmt := range functions {
		if s, ok := stmt.(*ast.FunctionDeclaration); ok {
			sb.WriteString(pg.GeneratePythonFunction(s))
		}
	}
	
	// Variables globales
	for _, stmt := range variables {
		if s, ok := stmt.(*ast.VariableDeclaration); ok {
			sb.WriteString(pg.GenerateVariableDeclaration(s))
		}
	}
	
	// Main execution
	if len(expressions) > 0 {
		sb.WriteString("\n# Main execution\n")
		for _, stmt := range expressions {
			if s, ok := stmt.(*ast.ExpressionStatement); ok {
				sb.WriteString(pg.GeneratePythonExpressionStatement(s))
			}
		}
	}

	return sb.String()
}

func (pg *PythonGenerator) GeneratePythonFunction(fd *ast.FunctionDeclaration) string {
	var sb strings.Builder
	
	sb.WriteString("def ")
	sb.WriteString(fd.Name)
	sb.WriteString("(")
	
	// Paramètres
	for i, param := range fd.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
	}
	
	sb.WriteString("):\n")
	
	// Corps de la fonction
	if len(fd.Body) == 0 {
		sb.WriteString("    pass\n")
	} else {
		for _, stmt := range fd.Body {
			sb.WriteString("    ")
			sb.WriteString(pg.GeneratePythonStatement(stmt))
		}
	}
	
	sb.WriteString("\n")
	return sb.String()
}

func (pg *PythonGenerator) GeneratePythonStatement(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		if s.Value != nil {
			return "return " + pg.GeneratePythonExpression(s.Value) + "\n"
		}
		return "return\n"
	case *ast.IfStatement:
		return pg.GeneratePythonIfStatement(s)
	case *ast.VariableDeclaration:
		return pg.GenerateVariableDeclaration(s)
	}
	return ""
}

func (pg *PythonGenerator) GeneratePythonIfStatement(is *ast.IfStatement) string {
	var sb strings.Builder
	sb.WriteString("if ")
	sb.WriteString(pg.GeneratePythonExpression(is.Condition))
	sb.WriteString(":\n")
	
	if blockStmt, ok := is.ThenBranch.(*ast.BlockStatement); ok {
		for _, stmt := range blockStmt.Statements {
			sb.WriteString("        ")
			sb.WriteString(pg.GeneratePythonStatement(stmt))
		}
	}
	
	if is.ElseBranch != nil {
		sb.WriteString("    else:\n")
		if blockStmt, ok := is.ElseBranch.(*ast.BlockStatement); ok {
			for _, stmt := range blockStmt.Statements {
				sb.WriteString("        ")
				sb.WriteString(pg.GeneratePythonStatement(stmt))
			}
		}
	}
	
	return sb.String()
}

func (pg *PythonGenerator) GeneratePythonExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return pg.GenerateStringLiteral(e)
	case *ast.NumberLiteral:
		return pg.GenerateNumberLiteral(e)
	case *ast.BooleanLiteral:
		if e.Value {
			return "True"
		}
		return "False"
	case *ast.TemplateLiteral:
		return pg.GenerateTemplateLiteral(e)
	case *ast.ArrayLiteral:
		return pg.GenerateArrayLiteral(e)
	case *ast.ObjectLiteral:
		return pg.GenerateObjectLiteral(e)
	case *ast.CallExpression:
		return pg.GenerateCallExpression(e)
	case *ast.IndexExpression:
		return pg.GenerateIndexExpression(e)
	case *ast.Identifier:
		return e.Value
	case *ast.InfixExpression:
		return pg.GeneratePythonExpression(e.Left) + " " + e.Operator + " " + pg.GeneratePythonExpression(e.Right)
	}
	return ""
}

func (pg *PythonGenerator) GeneratePythonExpressionStatement(es *ast.ExpressionStatement) string {
	// Convertir console.log en print
	if callExpr, ok := es.Expression.(*ast.CallExpression); ok {
		if dotExpr, ok := callExpr.Function.(*ast.DotExpression); ok {
			if ident, ok := dotExpr.Object.(*ast.Identifier); ok {
				if ident.Value == "console" && dotExpr.Property == "log" {
					var sb strings.Builder
					sb.WriteString("print(")
					for i, arg := range callExpr.Arguments {
						if i > 0 {
							sb.WriteString(", ")
						}
						sb.WriteString(pg.GeneratePythonExpression(arg))
					}
					sb.WriteString(")\n")
					return sb.String()
				}
			}
		}
	}
	
	return pg.GeneratePythonExpression(es.Expression) + "\n"
}

func (pg *PythonGenerator) GenerateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	// Python n'a pas de const, on peut utiliser un commentaire ou une convention
	if vd.IsConst {
		sb.WriteString("# Constant\n")
	}
	
	sb.WriteString(vd.Name)
	sb.WriteString(" = ")
	
	switch val := vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString(pg.GenerateStringLiteral(val))
	case *ast.NumberLiteral:
		sb.WriteString(pg.GenerateNumberLiteral(val))
	case *ast.BooleanLiteral:
		sb.WriteString(pg.GenerateBooleanLiteral(val))
	case *ast.ArrayLiteral:
		sb.WriteString(pg.GenerateArrayLiteral(val))
	case *ast.ObjectLiteral:
		sb.WriteString(pg.GenerateObjectLiteral(val))
	case *ast.TemplateLiteral:
		sb.WriteString(pg.GenerateTemplateLiteral(val))
	default:
		sb.WriteString(pg.GeneratePythonExpression(val))
	}
	
	sb.WriteString("\n")
	return sb.String()
}

func (pg *PythonGenerator) GenerateCallExpression(ce *ast.CallExpression) string {
	var sb strings.Builder
	sb.WriteString(pg.GeneratePythonExpression(ce.Function))
	sb.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(pg.GeneratePythonExpression(arg))
	}
	sb.WriteString(")")
	return sb.String()
}

func (pg *PythonGenerator) GenerateIndexExpression(ie *ast.IndexExpression) string {
	return pg.GeneratePythonExpression(ie.Left) + "[" + pg.GeneratePythonExpression(ie.Index) + "]"
}

func (pg *PythonGenerator) GenerateArrayLiteral(al *ast.ArrayLiteral) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, element := range al.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(pg.GeneratePythonExpression(element))
	}
	sb.WriteString("]")
	return sb.String()
}

func (pg *PythonGenerator) GenerateObjectLiteral(ol *ast.ObjectLiteral) string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, prop := range ol.Properties {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("\"")
		sb.WriteString(prop.Key)
		sb.WriteString("\": ")
		sb.WriteString(pg.GeneratePythonExpression(prop.Value))
	}
	sb.WriteString("}")
	return sb.String()
}

func (pg *PythonGenerator) GenerateStringLiteral(sl *ast.StringLiteral) string {
	return "\"" + sl.Value + "\""
}

func (pg *PythonGenerator) GenerateNumberLiteral(nl *ast.NumberLiteral) string {
	return nl.Value
}

func (pg *PythonGenerator) GenerateBooleanLiteral(bl *ast.BooleanLiteral) string {
	if bl.Value {
		return "True"
	}
	return "False"
}

func (pg *PythonGenerator) GenerateTemplateLiteral(tl *ast.TemplateLiteral) string {
	// En Python, convertir les template literals en f-strings
	if len(tl.Parts) > 0 {
		content := tl.Parts[0].TokenLiteral()
		// Pour l'instant, retourner comme string simple (plus tard on peut parser ${} pour interpolation)
		return "\"" + content + "\""
	}
	return "\"\""
}

// CSharpGenerator génère du code C#
type CSharpGenerator struct{}

func (csg *CSharpGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder
	
	sb.WriteString("using System;\n\n")
	sb.WriteString("namespace GeneratedCode\n{\n")
	sb.WriteString("    class Program\n    {\n")
	sb.WriteString("        static void Main(string[] args)\n        {\n")
	
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			sb.WriteString("            ")
			sb.WriteString(csg.GenerateVariableDeclaration(s))
		}
	}
	
	sb.WriteString("        }\n    }\n}\n")
	return sb.String()
}

func (csg *CSharpGenerator) GenerateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	// Déterminer le type C#
	switch vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString("string ")
	case *ast.NumberLiteral:
		sb.WriteString("int ") // Simplifié
	case *ast.BooleanLiteral:
		sb.WriteString("bool ")
	default:
		sb.WriteString("var ")
	}
	
	sb.WriteString(vd.Name)
	sb.WriteString(" = ")
	
	switch val := vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString(csg.GenerateStringLiteral(val))
	case *ast.NumberLiteral:
		sb.WriteString(csg.GenerateNumberLiteral(val))
	case *ast.BooleanLiteral:
		sb.WriteString(csg.GenerateBooleanLiteral(val))
	}
	
	sb.WriteString(";\n")
	return sb.String()
}

func (csg *CSharpGenerator) GenerateStringLiteral(sl *ast.StringLiteral) string {
	return "\"" + sl.Value + "\""
}

func (csg *CSharpGenerator) GenerateNumberLiteral(nl *ast.NumberLiteral) string {
	return nl.Value
}

func (csg *CSharpGenerator) GenerateBooleanLiteral(bl *ast.BooleanLiteral) string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// GoGenerator génère du code Go
type GoGenerator struct{}

func (gg *GoGenerator) Generate(statements []ast.Statement) string {
	var sb strings.Builder
	
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			sb.WriteString("    ")
			sb.WriteString(gg.GenerateVariableDeclaration(s))
		}
	}
	
	sb.WriteString("}\n")
	return sb.String()
}

func (gg *GoGenerator) GenerateVariableDeclaration(vd *ast.VariableDeclaration) string {
	var sb strings.Builder
	
	if vd.IsConst {
		sb.WriteString("const ")
	} else {
		sb.WriteString("var ")
	}
	
	sb.WriteString(vd.Name)
	sb.WriteString(" ")
	
	// Déterminer le type Go
	switch vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString("string")
	case *ast.NumberLiteral:
		sb.WriteString("int")
	case *ast.BooleanLiteral:
		sb.WriteString("bool")
	default:
		sb.WriteString("interface{}")
	}
	
	sb.WriteString(" = ")
	
	switch val := vd.Value.(type) {
	case *ast.StringLiteral:
		sb.WriteString(gg.GenerateStringLiteral(val))
	case *ast.NumberLiteral:
		sb.WriteString(gg.GenerateNumberLiteral(val))
	case *ast.BooleanLiteral:
		sb.WriteString(gg.GenerateBooleanLiteral(val))
	}
	
	sb.WriteString("\n")
	return sb.String()
}

func (gg *GoGenerator) GenerateStringLiteral(sl *ast.StringLiteral) string {
	return "\"" + sl.Value + "\""
}

func (gg *GoGenerator) GenerateNumberLiteral(nl *ast.NumberLiteral) string {
	return nl.Value
}

func (gg *GoGenerator) GenerateBooleanLiteral(bl *ast.BooleanLiteral) string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// Fonction de compatibilité pour l'ancien code
func GenerateJS(statements []ast.Statement) string {
	return Generate(statements, JavaScript)
}

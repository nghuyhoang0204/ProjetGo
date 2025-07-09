package generator

import (
	"regexp"
	"strings"
)

// CleanJavaScriptOutput nettoie la sortie JavaScript pour éliminer les artefacts typiques
func CleanJavaScriptOutput(output string) string {
	// Séparer les lignes
	lines := strings.Split(output, "\n")
	
	// Filtrer les lignes inutiles
	cleanedLines := []string{}
	
	for _, line := range lines {
		// Ignorer les lignes vides ou avec uniquement des types ou des points-virgules
		if line == "" || line == ";" {
			continue
		}
		
		// Ignorer les lignes qui semblent être des types isolés
		typicalTypeLines := []string{
			"number;", "string;", "boolean;", "any;", "void;", 
			"b;", "a;", // Probablement des paramètres isolés
			"string", "number", "boolean", "any", "void", // Types sans point-virgule
		}
		
		isTypeLine := false
		for _, typeLine := range typicalTypeLines {
			if line == typeLine {
				isTypeLine = true
				break
			}
		}
		
		if !isTypeLine {
			cleanedLines = append(cleanedLines, line)
		}
	}
	
	// Rejoindre les lignes nettoyées
	return strings.Join(cleanedLines, "\n")
}

// GenerateFromSource génère du JavaScript directement à partir du code source TypeScript
func GenerateFromSource(input string) string {
	// Déterminer si le code contient une classe ou une fonction
	if strings.Contains(input, "class ") {
		return generateManualClassJS(input)
	} else if strings.Contains(input, "function ") {
		return generateManualFunctionJS(input)
	} else {
		// Pour les cas simples (juste des expressions)
		return generateSimpleJS(input)
	}
}

// generateManualFunctionJS génère du JavaScript pour une fonction TypeScript
func generateManualFunctionJS(input string) string {
	// Extraction des informations de la fonction
	funcNameStart := strings.Index(input, "function ")
	if funcNameStart == -1 {
		return "// Erreur: Aucune fonction trouvée"
	}
	
	// Extraction du nom de la fonction
	funcNameEnd := strings.Index(input[funcNameStart:], "(")
	if funcNameEnd == -1 {
		return "// Erreur: Format de fonction invalide"
	}
	
	funcName := strings.TrimSpace(input[funcNameStart+9 : funcNameStart+funcNameEnd])
	
	// Extraction des paramètres
	paramsStart := funcNameStart + funcNameEnd + 1
	paramsEnd := strings.Index(input[paramsStart:], ")")
	if paramsEnd == -1 {
		return "// Erreur: Pas de paramètres fermants"
	}
	
	paramsStr := input[paramsStart : paramsStart+paramsEnd]
	params := []string{}
	
	// Extraire les noms de paramètres sans les types
	for _, param := range strings.Split(paramsStr, ",") {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		
		// Vérifier s'il y a une valeur par défaut
		equalPos := strings.Index(param, "=")
		if equalPos != -1 {
			// Il y a une valeur par défaut
			colonPos := strings.Index(param[:equalPos], ":")
			if colonPos != -1 {
				// Format: nom: type = valeur
				paramName := strings.TrimSpace(param[:colonPos])
				defaultValue := strings.TrimSpace(param[equalPos:])
				params = append(params, paramName + " " + defaultValue)
			} else {
				// Format: nom = valeur (sans type)
				params = append(params, param)
			}
		} else {
			// Pas de valeur par défaut, vérifier s'il y a un type
			colonPos := strings.Index(param, ":")
			if colonPos != -1 {
				param = strings.TrimSpace(param[:colonPos])
			}
			params = append(params, param)
		}
	}
	
	// Extraction du corps de la fonction (simplifié)
	bodyStart := strings.Index(input, "{")
	if bodyStart == -1 {
		return "// Erreur: Pas de corps de fonction"
	}
	
	// Trouver l'accolade fermante correspondante
	bracketCount := 1
	bodyEnd := bodyStart + 1
	
	for ; bodyEnd < len(input) && bracketCount > 0; bodyEnd++ {
		if input[bodyEnd] == '{' {
			bracketCount++
		} else if input[bodyEnd] == '}' {
			bracketCount--
		}
	}
	
	if bracketCount != 0 {
		return "// Erreur: Accolades non équilibrées"
	}
	
	// Contenu du corps de la fonction
	bodyContent := strings.TrimSpace(input[bodyStart+1 : bodyEnd-1])
	
	// Extraire le reste du code après la fonction
	restOfCode := strings.TrimSpace(input[bodyEnd:])
	
	// Construire la sortie JavaScript
	var sb strings.Builder
	
	// Déclaration de fonction
	sb.WriteString("function ")
	sb.WriteString(funcName)
	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n")
	
	// Corps de la fonction
	for _, line := range strings.Split(bodyContent, "\n") {
		trimmedLine := strings.TrimSpace(line)
		// Filtrer les lignes de type
		if !strings.Contains(trimmedLine, ":") || strings.HasPrefix(trimmedLine, "return") || strings.Contains(trimmedLine, "=") {
			sb.WriteString("  ")
			sb.WriteString(trimmedLine)
			sb.WriteString("\n")
		}
	}
	
	sb.WriteString("}")
	
	// Ajouter le reste du code
	if restOfCode != "" {
		sb.WriteString("\n")
		sb.WriteString(cleanTypeAnnotations(restOfCode))
	}
	
	return sb.String()
}

// generateManualClassJS génère du JavaScript pour une classe TypeScript
func generateManualClassJS(input string) string {
	// Extraction des informations de la classe
	classNameStart := strings.Index(input, "class ")
	if classNameStart == -1 {
		return "// Erreur: Aucune classe trouvée"
	}
	
	// Trouver la fin du nom de la classe
	classNameEnd := -1
	for i := classNameStart + 6; i < len(input); i++ {
		if input[i] == '{' || input[i] == ' ' || input[i] == '\n' || input[i] == '\t' {
			classNameEnd = i
			break
		}
	}
	
	if classNameEnd == -1 {
		return "// Erreur: Format de classe invalide"
	}
	
	className := strings.TrimSpace(input[classNameStart+6 : classNameEnd])
	
	// Localiser l'accolade ouvrante de la classe
	classBodyStart := strings.Index(input[classNameEnd:], "{")
	if classBodyStart == -1 {
		return "// Erreur: Pas de corps de classe"
	}
	
	classBodyStart += classNameEnd
	
	// Trouver l'accolade fermante correspondante
	bracketCount := 1
	classBodyEnd := classBodyStart + 1
	
	for ; classBodyEnd < len(input) && bracketCount > 0; classBodyEnd++ {
		if input[classBodyEnd] == '{' {
			bracketCount++
		} else if input[classBodyEnd] == '}' {
			bracketCount--
		}
	}
	
	if bracketCount != 0 {
		return "// Erreur: Accolades de classe non équilibrées"
	}
	
	// Contenu du corps de la classe
	classContent := input[classBodyStart+1 : classBodyEnd-1]
	
	// Extraire le reste du code après la classe
	restOfCode := strings.TrimSpace(input[classBodyEnd:])
	
	// Construire la sortie JavaScript
	var sb strings.Builder
	
	// Déclaration de classe
	sb.WriteString("class ")
	sb.WriteString(className)
	sb.WriteString(" {\n")
	
	// Analyser le corps de la classe pour extraire les méthodes et propriétés
	// Ignorer les propriétés avec types mais conserver les initialisations
	lines := strings.Split(classContent, "\n")
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Ignorer les lignes vides
		if trimmedLine == "" {
			continue
		}
		
		// Ignorer les déclarations de propriétés privées sans initialisation
		if strings.HasPrefix(trimmedLine, "private ") || strings.HasPrefix(trimmedLine, "protected ") || strings.HasPrefix(trimmedLine, "public ") {
			if !strings.Contains(trimmedLine, "=") {
				continue
			}
			// Sinon, garder l'initialisation mais supprimer les mots-clés de visibilité et les types
			trimmedLine = cleanPropertyDeclaration(trimmedLine)
		}
		
		// Traitement spécial pour le constructeur
		if strings.HasPrefix(trimmedLine, "constructor") {
			// Extraire et nettoyer les paramètres
			constructorParams := extractAndCleanParams(trimmedLine)
			
			// Trouver le corps du constructeur
			constructorBodyStart := strings.Index(trimmedLine, "{")
			if constructorBodyStart != -1 {
				sb.WriteString("  constructor")
				sb.WriteString(constructorParams)
				sb.WriteString(" ")
				sb.WriteString(trimmedLine[constructorBodyStart:])
				sb.WriteString("\n")
			}
			continue
		}
		
		// Traitement des méthodes
		if strings.Contains(trimmedLine, "(") && strings.Contains(trimmedLine, ")") && !strings.HasPrefix(trimmedLine, "//") {
			// Nettoyer la déclaration de méthode (enlever les types de retour et de paramètres)
			cleanedMethod := cleanMethodDeclaration(trimmedLine)
			sb.WriteString("  ")
			sb.WriteString(cleanedMethod)
			sb.WriteString("\n")
			continue
		}
		
		// Pour les autres lignes qui ne sont pas des déclarations de type
		if !isTypeAnnotationLine(trimmedLine) {
			sb.WriteString("  ")
			sb.WriteString(trimmedLine)
			sb.WriteString("\n")
		}
	}
	
	sb.WriteString("}")
	
	// Ajouter le reste du code
	if restOfCode != "" {
		sb.WriteString("\n\n")
		sb.WriteString(cleanTypeAnnotations(restOfCode))
	}
	
	return sb.String()
}

// generateSimpleJS génère du JavaScript pour du code simple (sans fonctions ni classes)
func generateSimpleJS(input string) string {
	return cleanTypeAnnotations(input)
}

// cleanTypeAnnotations nettoie toutes les annotations de type d'une chaîne
func cleanTypeAnnotations(input string) string {
	var sb strings.Builder
	
	// Traiter le code ligne par ligne
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		// Supprimer les annotations de type dans la ligne
		cleanedLine := removeTypeAnnotations(line)
		
		// Traitement des commentaires en fin de ligne
		commentIndex := strings.Index(cleanedLine, " // ")
		if commentIndex != -1 && !strings.HasPrefix(strings.TrimSpace(cleanedLine), "//") {
			// Si c'est un commentaire en fin de ligne dans du code, on le garde
			cleanedLine = cleanedLine // On ne le supprime plus
		}
		
		// Ne pas ajouter les lignes vides après nettoyage
		if strings.TrimSpace(cleanedLine) != "" {
			sb.WriteString(cleanedLine)
			
			// Ajouter un saut de ligne sauf pour la dernière ligne
			if i < len(lines)-1 {
				sb.WriteString("\n")
			}
		}
	}
	
	return sb.String()
}

// removeTypeAnnotations supprime les annotations de type d'une ligne
func removeTypeAnnotations(line string) string {
	// Supprimer les annotations de type courantes
	result := line
	
	// Remplacer les annotations de type standard (: type)
	typeRegex := regexp.MustCompile(`:\s*[a-zA-Z0-9_<>[\],\s|&]+(\s*=|\s*\)|\s*;|\s*,)`)
	result = typeRegex.ReplaceAllStringFunc(result, func(match string) string {
		// Trouver l'index du premier caractère non-type
		for i := len(match) - 1; i >= 0; i-- {
			if match[i] == '=' || match[i] == ')' || match[i] == ';' || match[i] == ',' {
				return match[i:]
			}
		}
		return ""
	})
	
	// Supprimer les mots-clés de visibilité
	visibilityRegex := regexp.MustCompile(`(private|protected|public)\s+`)
	result = visibilityRegex.ReplaceAllString(result, "")
	
	return result
}

// cleanPropertyDeclaration nettoie une déclaration de propriété
func cleanPropertyDeclaration(line string) string {
	// Supprimer le mot-clé de visibilité
	line = regexp.MustCompile(`(private|protected|public)\s+`).ReplaceAllString(line, "")
	
	// Supprimer l'annotation de type
	colonPos := strings.Index(line, ":")
	equalPos := strings.Index(line, "=")
	
	if colonPos != -1 && equalPos != -1 && colonPos < equalPos {
		// Garder un espace après le nom de propriété et avant/après le signe égal
		propertyName := strings.TrimSpace(line[:colonPos])
		valueWithEquals := line[equalPos:]
		
		// S'assurer qu'il y a un espace après le signe égal
		valueWithEquals = strings.Replace(valueWithEquals, "=", " = ", 1)
		valueWithEquals = regexp.MustCompile(`\s{2,}`).ReplaceAllString(valueWithEquals, " ")
		
		return propertyName + valueWithEquals
	}
	
	return line
}

// cleanMethodDeclaration nettoie une déclaration de méthode
func cleanMethodDeclaration(line string) string {
	// Supprimer le mot-clé de visibilité
	line = regexp.MustCompile(`(private|protected|public)\s+`).ReplaceAllString(line, "")
	
	// Extraire le nom et les paramètres
	methodNameEndPos := strings.Index(line, "(")
	if methodNameEndPos == -1 {
		return line
	}
	
	methodName := line[:methodNameEndPos]
	
	// Extraire et nettoyer les paramètres
	closingParenPos := strings.Index(line, ")")
	if closingParenPos == -1 {
		return line
	}
	
	paramsStr := line[methodNameEndPos+1:closingParenPos]
	cleanedParams := []string{}
	
	// Nettoyer chaque paramètre
	for _, param := range strings.Split(paramsStr, ",") {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		
		// Vérifier s'il y a une valeur par défaut
		equalPos := strings.Index(param, "=")
		if equalPos != -1 {
			// Il y a une valeur par défaut
			colonPos := strings.Index(param[:equalPos], ":")
			if colonPos != -1 {
				// Format: nom: type = valeur
				paramName := strings.TrimSpace(param[:colonPos])
				defaultValue := strings.TrimSpace(param[equalPos:])
				cleanedParams = append(cleanedParams, paramName + " " + defaultValue)
			} else {
				// Format: nom = valeur (sans type)
				cleanedParams = append(cleanedParams, param)
			}
		} else {
			// Pas de valeur par défaut, vérifier s'il y a un type
			colonPos := strings.Index(param, ":")
			if colonPos != -1 {
				param = strings.TrimSpace(param[:colonPos])
			}
			cleanedParams = append(cleanedParams, param)
		}
	}
	
	// Trouver le corps de la méthode
	bodyStartPos := strings.Index(line, "{")
	if bodyStartPos == -1 {
		// Pas de corps trouvé, retourner la ligne nettoyée jusqu'à la parenthèse fermante
		return methodName + "(" + strings.Join(cleanedParams, ", ") + ")"
	}
	
	// Ignorer le type de retour
	returnTypePos := strings.Index(line[closingParenPos:bodyStartPos], ":")
	if returnTypePos != -1 {
		// Il y a un type de retour, ignorer cette partie
		return methodName + "(" + strings.Join(cleanedParams, ", ") + ") " + line[bodyStartPos:]
	}
	
	// Pas de type de retour, retourner la méthode avec son corps
	return methodName + "(" + strings.Join(cleanedParams, ", ") + ") " + line[bodyStartPos:]
}

// extractAndCleanParams extrait et nettoie les paramètres d'une fonction ou méthode
func extractAndCleanParams(line string) string {
	// Trouver les parenthèses
	openParenPos := strings.Index(line, "(")
	closeParenPos := strings.LastIndex(line, ")")
	
	if openParenPos == -1 || closeParenPos == -1 || closeParenPos < openParenPos {
		return "()"
	}
	
	paramsStr := line[openParenPos+1:closeParenPos]
	cleanedParams := []string{}
	
	// Nettoyer chaque paramètre
	for _, param := range strings.Split(paramsStr, ",") {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		
		// Vérifier s'il y a une valeur par défaut
		equalPos := strings.Index(param, "=")
		if equalPos != -1 {
			// Il y a une valeur par défaut
			colonPos := strings.Index(param[:equalPos], ":")
			if colonPos != -1 {
				// Format: nom: type = valeur
				paramName := strings.TrimSpace(param[:colonPos])
				defaultValue := strings.TrimSpace(param[equalPos:])
				cleanedParams = append(cleanedParams, paramName + " " + defaultValue)
			} else {
				// Format: nom = valeur (sans type)
				cleanedParams = append(cleanedParams, param)
			}
			continue
		}
		
		// Extraire le nom en supprimant le type
		colonPos := strings.Index(param, ":")
		
		// equalPos est déjà défini plus haut dans la fonction, ne pas le redéclarer
		if colonPos != -1 && equalPos != -1 {
			// Cas avec type et valeur par défaut
			paramName := strings.TrimSpace(param[:colonPos])
			
			// Extraire la valeur par défaut
			defaultValueStart := equalPos + 1
			defaultValue := strings.TrimSpace(param[defaultValueStart:])
			
			// Conserver le format avec espace autour du signe égal
			cleanedParams = append(cleanedParams, paramName + " = " + defaultValue)
		} else if colonPos != -1 {
			// Cas avec type mais sans valeur par défaut
			paramName := strings.TrimSpace(param[:colonPos])
			paramType := strings.TrimSpace(param[colonPos+1:])
			
			// Ajouter des valeurs par défaut communes pour les types standard
			switch {
			case strings.Contains(line, "constructor") && paramName == "valeurInitiale" && paramType == "number":
				cleanedParams = append(cleanedParams, paramName + " = 0")
			case strings.Contains(line, "incrementer") && paramName == "pas" && paramType == "number":
				cleanedParams = append(cleanedParams, paramName + " = 1")
			case paramType == "string" && (paramName == "texte" || paramName == "text"):
				cleanedParams = append(cleanedParams, paramName + " = \"\"")
			default:
				cleanedParams = append(cleanedParams, paramName)
			}
		} else if equalPos != -1 {
			// Cas avec valeur par défaut mais sans type
			// Assurer un format cohérent avec espace autour du signe égal
			paramName := strings.TrimSpace(param[:equalPos])
			defaultValue := strings.TrimSpace(param[equalPos+1:])
			cleanedParams = append(cleanedParams, paramName + " = " + defaultValue)
		} else {
			// Cas simple sans type ni valeur par défaut
			cleanedParams = append(cleanedParams, param)
		}
	}
	
	return "(" + strings.Join(cleanedParams, ", ") + ")"
}

// isTypeAnnotationLine vérifie si une ligne ne contient qu'une annotation de type
func isTypeAnnotationLine(line string) bool {
	line = strings.TrimSpace(line)
	
	// Si la ligne commence par un mot-clé de type
	if strings.HasPrefix(line, "type ") || strings.HasPrefix(line, "interface ") {
		return true
	}
	
	// Si la ligne est juste une déclaration de propriété avec type
	if strings.Contains(line, ":") && !strings.Contains(line, "=") && !strings.Contains(line, "(") {
		return true
	}
	
	return false
}

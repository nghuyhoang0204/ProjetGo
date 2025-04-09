package parser

import (
	abstract "ProjetGo/ast"
	"ProjetGo/lexer"
)

// Parse converts a list of tokens into an AST.
func Parse(tokens []lexer.Token) abstract.AST {
	var astNodes []abstract.Node

	// A simple parser that only handles variable declarations for now
	for i := 0; i < len(tokens)-3; i++ { // Ensure there are enough tokens to parse
		if tokens[i].Value == "let" { // Check for variable declaration
			varName := tokens[i+1].Value  // Get the variable name
			if tokens[i+2].Value == ":" { // Skip the type annotation
				i += 2 // Move past the ":" and the type annotation
			}
			if tokens[i+2].Value == "=" { // Check for "="
				varValue := tokens[i+3].Value // Get the variable value
				astNodes = append(astNodes, &abstract.VariableDeclaration{Name: varName, Value: varValue})
				i += 3 // Skip past the declaration
			}
		}
	}

	return abstract.AST{Nodes: astNodes}
}

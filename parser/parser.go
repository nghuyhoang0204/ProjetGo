package parser

import (
	"ProjetGo/ast"
	"ProjetGo/lexer"
)

// Parse converts a list of tokens into an AST.
func Parse(tokens []lexer.Token) ast.AST {
	var astNodes []ast.Node

	// A simple parser that only handles variable declarations for now
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Value == "let" {
			varName := tokens[i+1].Value
			varValue := tokens[i+3].Value
			astNodes = append(astNodes, &ast.VariableDeclaration{Name: varName, Value: varValue})
			i += 4 // Skip past the declaration
		}
	}

	return ast.AST{Nodes: astNodes}
}

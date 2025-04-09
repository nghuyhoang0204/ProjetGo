package codegen

import (
	"ProjetGo/ast"
	"strings"
)

// GenerateCode generates JavaScript code from the AST.
func GenerateCode(ast ast.AST) string {
	var jsCode []string

	// Iterate over the AST nodes and generate JavaScript code
	for _, node := range ast.Nodes {
		switch n := node.(type) {
		case *ast.VariableDeclaration: // Correct usage of the pointer type
			jsCode = append(jsCode, "let "+n.Name+" = "+n.Value+";")
		}
	}

	return strings.Join(jsCode, "\n")
}

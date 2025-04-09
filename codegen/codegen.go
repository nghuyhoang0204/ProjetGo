package codegen

import (
	abstract "ProjetGo/ast"
	"fmt"
	"strings"
)

func GenerateCode(ast abstract.AST) string {
	var jsCode []string

	for _, node := range ast.Nodes {
		switch n := node.(type) {
		case *abstract.VariableDeclaration:
			fmt.Println("Generating code for variable:", n.Name)
			jsCode = append(jsCode, "let "+n.Name+" = "+n.Value+";")
		}
	}

	fmt.Println("Generated JavaScript:", jsCode)
	return strings.Join(jsCode, "\n")
}

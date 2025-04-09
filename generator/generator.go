package generator

import (
	"strings"
	"ProjetGo/ast"
)

func GenerateJS(statements []ast.Statement) string {
	var sb strings.Builder

	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			if s.IsConst {
				sb.WriteString("const ")
			} else {
				sb.WriteString("let ")
			}
			sb.WriteString(s.Name)
			sb.WriteString(" = ")

			switch val := s.Value.(type) {
			case *ast.StringLiteral:
				sb.WriteString("\"" + val.Value + "\"")
			case *ast.NumberLiteral:
				sb.WriteString(val.Value)
			}

			sb.WriteString(";\n")
		}
	}

	return sb.String()
}

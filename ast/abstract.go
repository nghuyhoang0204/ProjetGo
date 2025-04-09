package abstract

// Node represents an AST node interface.
type Node interface {
	Type() string
}

// VariableDeclaration represents a variable declaration in TypeScript.
type VariableDeclaration struct {
	Name  string
	Value string
}

func (v *VariableDeclaration) Type() string {
	return "VariableDeclaration"
}

// AST holds all the nodes of the parsed TypeScript code.
type AST struct {
	Nodes []Node
}

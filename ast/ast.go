package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// VariableDeclaration est une instruction de type const/let/var
type VariableDeclaration struct {
	IsConst bool
	Name    string
	Type    string
	Value   Expression
}

func (vd *VariableDeclaration) statementNode()      {}
func (vd *VariableDeclaration) TokenLiteral() string { return "var" }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string { return s.Value }

type NumberLiteral struct {
	Value string
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) TokenLiteral() string { return n.Value }

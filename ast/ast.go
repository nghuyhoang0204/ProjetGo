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

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) TokenLiteral() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// TypeAlias pour les alias de types comme type TaskStatus = 'pending' | 'in_progress' | 'done'
type TypeAlias struct {
	Name string
	Type string
}

func (ta *TypeAlias) statementNode() {}
func (ta *TypeAlias) TokenLiteral() string { return "type" }

// Interface pour les interfaces TypeScript
type Interface struct {
	Name   string
	Fields []InterfaceField
}

func (i *Interface) statementNode() {}
func (i *Interface) TokenLiteral() string { return "interface" }

type InterfaceField struct {
	Name string
	Type string
}

// ClassDeclaration pour les classes
type ClassDeclaration struct {
	Name    string
	Fields  []ClassField
	Methods []ClassMethod
}

func (cd *ClassDeclaration) statementNode() {}
func (cd *ClassDeclaration) TokenLiteral() string { return "class" }

type ClassField struct {
	Name       string
	Type       string
	IsPrivate  bool
	IsStatic   bool
	HasDefault bool
	Default    Expression
}

type ClassMethod struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	IsAsync    bool
	IsPrivate  bool
	Body       []Statement
}

type Parameter struct {
	Name string
	Type string
}

// FunctionDeclaration pour les fonctions
type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	IsAsync    bool
	Body       []Statement
}

func (fd *FunctionDeclaration) statementNode() {}
func (fd *FunctionDeclaration) TokenLiteral() string { return "function" }

// ArrayLiteral pour les tableaux
type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string { return "[" }

// ObjectLiteral pour les objets
type ObjectLiteral struct {
	Properties []ObjectProperty
}

func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) TokenLiteral() string { return "{" }

type ObjectProperty struct {
	Key   string
	Value Expression
}

// Statements pour le contrôle de flux
type IfStatement struct {
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) TokenLiteral() string { return "if" }

type ForStatement struct {
	Init      Statement
	Condition Expression
	Update    Statement
	Body      Statement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) TokenLiteral() string { return "for" }

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string { return "while" }

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return "{" }

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Expression.TokenLiteral() }

type ReturnStatement struct {
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return "return" }

// Expressions
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Value }

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return "(" }

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Operator }

// Template literals pour les backticks
type TemplateLiteral struct {
	Parts []Expression // Alternance de texte et expressions
}

func (tl *TemplateLiteral) expressionNode() {}
func (tl *TemplateLiteral) TokenLiteral() string { return "`" }

// Index access pour arr[0] ou obj.prop
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string { return "[" }

type DotExpression struct {
	Object   Expression
	Property string
}

func (de *DotExpression) expressionNode() {}
func (de *DotExpression) TokenLiteral() string { return "." }

// Assignment pour depart--
type AssignmentExpression struct {
	Left     Expression
	Operator string // =, +=, -=, ++, --
	Right    Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Operator }

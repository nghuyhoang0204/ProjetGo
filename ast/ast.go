package ast

// Node est l'interface de base pour tous les nœuds de l'AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement représente les instructions (déclarations, conditions, boucles, etc.)
type Statement interface {
	Node
	statementNode()
}

// Expression représente les expressions (valeurs, opérations, appels de fonction, etc.)
type Expression interface {
	Node
	expressionNode()
}

// =============================================================================
// STATEMENTS (Instructions)
// =============================================================================

// Program est le nœud racine qui contient tous les statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

// VariableDeclaration représente: let/const/var name: type = value
type VariableDeclaration struct {
	Token     string // "let", "const", "var"
	Name      string
	Type      string // optionnel en TypeScript
	Value     Expression
	IsConst   bool
	IsExported bool // pour export
}

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token }
func (vd *VariableDeclaration) String() string {
	out := vd.Token + " " + vd.Name
	if vd.Type != "" {
		out += ": " + vd.Type
	}
	if vd.Value != nil {
		out += " = " + vd.Value.String()
	}
	return out + ";"
}

// FunctionDeclaration représente: function name(params): returnType { body }
type FunctionDeclaration struct {
	Token      string // "function"
	Name       string
	Parameters []Parameter
	ReturnType string
	Body       *BlockStatement
	IsAsync    bool
	IsExported bool
}

func (fd *FunctionDeclaration) statementNode()       {}
func (fd *FunctionDeclaration) TokenLiteral() string { return fd.Token }
func (fd *FunctionDeclaration) String() string {
	out := "function " + fd.Name + "("
	for i, p := range fd.Parameters {
		if i > 0 {
			out += ", "
		}
		out += p.String()
	}
	out += ")"
	if fd.ReturnType != "" {
		out += ": " + fd.ReturnType
	}
	out += " " + fd.Body.String()
	return out
}

// Parameter représente un paramètre de fonction
type Parameter struct {
	Name         string
	Type         string
	DefaultValue Expression
	IsOptional   bool
}

func (p *Parameter) String() string {
	out := p.Name
	if p.IsOptional {
		out += "?"
	}
	if p.Type != "" {
		out += ": " + p.Type
	}
	if p.DefaultValue != nil {
		out += " = " + p.DefaultValue.String()
	}
	return out
}

// IfStatement représente: if (condition) { then } else { else }
type IfStatement struct {
	Token     string // "if"
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement // optionnel
}

func (ifs *IfStatement) statementNode()       {}
func (ifs *IfStatement) TokenLiteral() string { return ifs.Token }
func (ifs *IfStatement) String() string {
	out := "if (" + ifs.Condition.String() + ") " + ifs.ThenBranch.String()
	if ifs.ElseBranch != nil {
		out += " else " + ifs.ElseBranch.String()
	}
	return out
}

// ForStatement représente: for (init; condition; update) { body }
type ForStatement struct {
	Token     string // "for"
	Init      Statement
	Condition Expression
	Update    Statement
	Body      Statement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token }
func (fs *ForStatement) String() string {
	out := "for ("
	if fs.Init != nil {
		out += fs.Init.String()
	}
	out += "; "
	if fs.Condition != nil {
		out += fs.Condition.String()
	}
	out += "; "
	if fs.Update != nil {
		out += fs.Update.String()
	}
	out += ") " + fs.Body.String()
	return out
}

// WhileStatement représente: while (condition) { body }
type WhileStatement struct {
	Token     string // "while"
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token }
func (ws *WhileStatement) String() string {
	return "while (" + ws.Condition.String() + ") " + ws.Body.String()
}

// BlockStatement représente: { statements... }
type BlockStatement struct {
	Token      string // "{"
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token }
func (bs *BlockStatement) String() string {
	out := "{ "
	for _, s := range bs.Statements {
		out += s.String() + " "
	}
	out += "}"
	return out
}

// ReturnStatement représente: return expression;
type ReturnStatement struct {
	Token string // "return"
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token }
func (rs *ReturnStatement) String() string {
	out := rs.Token
	if rs.Value != nil {
		out += " " + rs.Value.String()
	}
	return out + ";"
}

// ExpressionStatement représente une expression utilisée comme statement
type ExpressionStatement struct {
	Token      string
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token }
func (es *ExpressionStatement) String() string {
	return es.Expression.String() + ";"
}

// AssignmentStatement représente: variable = value
type AssignmentStatement struct {
	Token string
	Name  string
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token }
func (as *AssignmentStatement) String() string {
	return as.Name + " = " + as.Value.String() + ";"
}

// =============================================================================
// EXPRESSIONS (Valeurs et opérations)
// =============================================================================

// Identifier représente un nom de variable/fonction
type Identifier struct {
	Token string
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// StringLiteral représente: "string" ou 'string'
type StringLiteral struct {
	Token string
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// NumberLiteral représente: 123, 45.67
type NumberLiteral struct {
	Token string
	Value string
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token }
func (nl *NumberLiteral) String() string       { return nl.Value }

// BooleanLiteral représente: true, false
type BooleanLiteral struct {
	Token string
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token }
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// ArrayLiteral représente: [1, 2, 3]
type ArrayLiteral struct {
	Token    string // "["
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token }
func (al *ArrayLiteral) String() string {
	out := "["
	for i, e := range al.Elements {
		if i > 0 {
			out += ", "
		}
		out += e.String()
	}
	out += "]"
	return out
}

// ObjectLiteral représente: { key: value, ... }
type ObjectLiteral struct {
	Token      string // "{"
	Properties []ObjectProperty
}

type ObjectProperty struct {
	Key   string
	Value Expression
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) TokenLiteral() string { return ol.Token }
func (ol *ObjectLiteral) String() string {
	out := "{ "
	for i, p := range ol.Properties {
		if i > 0 {
			out += ", "
		}
		out += p.Key + ": " + p.Value.String()
	}
	out += " }"
	return out
}

// CallExpression représente: func(args...)
type CallExpression struct {
	Token     string // "("
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token }
func (ce *CallExpression) String() string {
	out := ce.Function.String() + "("
	for i, a := range ce.Arguments {
		if i > 0 {
			out += ", "
		}
		out += a.String()
	}
	out += ")"
	return out
}

// InfixExpression représente: left operator right
type InfixExpression struct {
	Token    string
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// MemberExpression représente: object.property ou object[key]
type MemberExpression struct {
	Token      string
	Object     Expression
	Property   Expression
	Computed   bool // true pour obj[key], false pour obj.prop
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token }
func (me *MemberExpression) String() string {
	if me.Computed {
		return me.Object.String() + "[" + me.Property.String() + "]"
	}
	return me.Object.String() + "." + me.Property.String()
}

// TemplateLiteral représente: `template ${expr} string`
type TemplateLiteral struct {
	Token string
	Parts []Expression // alternance string/expression
}

func (tl *TemplateLiteral) expressionNode()      {}
func (tl *TemplateLiteral) TokenLiteral() string { return tl.Token }
func (tl *TemplateLiteral) String() string {
	out := "`"
	for _, part := range tl.Parts {
		if _, ok := part.(*StringLiteral); ok {
			out += part.(*StringLiteral).Value
		} else {
			out += "${" + part.String() + "}"
		}
	}
	out += "`"
	return out
}

// TypeAlias représente une déclaration de type: type Name = Type
type TypeAlias struct {
	Token string // "type"
	Name  string
	Type  string
}

func (ta *TypeAlias) statementNode()  {}
func (ta *TypeAlias) TokenLiteral() string { return ta.Token }
func (ta *TypeAlias) String() string {
	var out string
	out += ta.TokenLiteral() + " " + ta.Name + " = " + ta.Type + ";"
	return out
}

// Interface représente une déclaration d'interface: interface Name { ... }
type Interface struct {
	Token      string // "interface"
	Name       string
	Properties []InterfaceProperty
}

type InterfaceProperty struct {
	Name       string
	Type       string
	IsOptional bool
}

func (i *Interface) statementNode()  {}
func (i *Interface) TokenLiteral() string { return i.Token }
func (i *Interface) String() string {
	var out string
	out += i.TokenLiteral() + " " + i.Name + " {\n"
	for _, prop := range i.Properties {
		out += "  " + prop.Name
		if prop.IsOptional {
			out += "?"
		}
		out += ": " + prop.Type + ";\n"
	}
	out += "}"
	return out
}

// ClassDeclaration représente une classe: class Name { ... }
type ClassDeclaration struct {
	Token      string // "class"
	Name       string
	SuperClass string // Optionnel, pour l'héritage
	Properties []ClassProperty
	Methods    []ClassMethod
}

type ClassProperty struct {
	Name       string
	Type       string
	IsPrivate  bool
	IsReadonly bool
	Value      Expression // Valeur initiale (optionnelle)
}

type ClassMethod struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	Body       *BlockStatement
	IsStatic   bool
	IsPrivate  bool
}

func (cd *ClassDeclaration) statementNode()  {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token }
func (cd *ClassDeclaration) String() string {
	var out string
	out += cd.TokenLiteral() + " " + cd.Name
	if cd.SuperClass != "" {
		out += " extends " + cd.SuperClass
	}
	out += " {\n"
	
	// Properties
	for _, prop := range cd.Properties {
		if prop.IsPrivate {
			out += "  private "
		}
		if prop.IsReadonly {
			out += "readonly "
		}
		out += prop.Name + ": " + prop.Type
		if prop.Value != nil {
			out += " = " + prop.Value.String()
		}
		out += ";\n"
	}
	
	// Methods
	for _, method := range cd.Methods {
		if method.IsStatic {
			out += "  static "
		}
		if method.IsPrivate {
			out += "private "
		}
		out += method.Name + "("
		
		// Parameters
		for i, param := range method.Parameters {
			if i > 0 {
				out += ", "
			}
			out += param.Name + ": " + param.Type
		}
		
		out += ")"
		if method.ReturnType != "" {
			out += ": " + method.ReturnType
		}
		out += " " + method.Body.String() + "\n"
	}
	
	out += "}"
	return out
}

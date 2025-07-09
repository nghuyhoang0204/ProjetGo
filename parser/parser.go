package parser

import (
	"ProjetGo/lexer"
	"ProjetGo/ast"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Charge deux tokens pour le lookahead
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
func (p *Parser) ParseStatement() ast.Statement {
	// Ignorer les commentaires
	for p.curToken.Type == lexer.COMMENT {
		p.nextToken()
	}
	
	switch p.curToken.Literal {
	case "let", "const", "var":
		return p.parseVariableDeclaration()
	case "function":
		return p.parseFunction()
	case "if":
		return p.parseIfStatement()
	case "for":
		return p.parseForStatement()
	case "while":
		return p.parseWhileStatement()
	case "return":
		return p.parseReturnStatement()
	case "type":
		return p.parseTypeAlias()
	case "interface":
		return p.parseInterface()
	case "class":
		return p.parseClass()
	default:
		// Essayer de parser comme expression statement
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfStatement() ast.Statement {
	p.nextToken() // passer 'if'
	
	if p.curToken.Type != lexer.LPAREN {
		return nil
	}
	p.nextToken() // passer '('
	
	condition := p.parseExpression()
	
	if p.curToken.Type != lexer.RPAREN {
		return nil
	}
	p.nextToken() // passer ')'
	
	thenBranch := p.parseBlockStatement()
	
	var elseBranch ast.Statement
	if p.curToken.Literal == "else" {
		p.nextToken()
		elseBranch = p.parseBlockStatement()
	}
	
	return &ast.IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	p.nextToken() // passer 'for'
	
	if p.curToken.Type != lexer.LPAREN {
		return nil
	}
	p.nextToken() // passer '('
	
	init := p.ParseStatement()
	condition := p.parseExpression()
	
	if p.curToken.Type != lexer.SEMICOLON {
		return nil
	}
	p.nextToken() // passer ';'
	
	update := p.ParseStatement()
	
	if p.curToken.Type != lexer.RPAREN {
		return nil
	}
	p.nextToken() // passer ')'
	
	body := p.parseBlockStatement()
	
	return &ast.ForStatement{
		Init:      init,
		Condition: condition,
		Update:    update,
		Body:      body,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	p.nextToken() // passer 'while'
	
	if p.curToken.Type != lexer.LPAREN {
		return nil
	}
	p.nextToken() // passer '('
	
	condition := p.parseExpression()
	
	if p.curToken.Type != lexer.RPAREN {
		return nil
	}
	p.nextToken() // passer ')'
	
	body := p.parseBlockStatement()
	
	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	p.nextToken() // passer 'return'
	
	var value ast.Expression
	if p.curToken.Type != lexer.SEMICOLON {
		value = p.parseExpression()
	}
	
	return &ast.ReturnStatement{Value: value}
}

func (p *Parser) parseBlockStatement() ast.Statement {
	if p.curToken.Type != lexer.LBRACE {
		return nil
	}
	p.nextToken() // passer '{'
	
	var statements []ast.Statement
	
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		} else {
			// Si ParseStatement retourne nil, avancer pour éviter la boucle infinie
			p.nextToken()
		}
		// Ne pas faire p.nextToken() ici car ParseStatement() gère déjà l'avancement
	}
	
	if p.curToken.Type == lexer.RBRACE {
		p.nextToken() // passer '}'
	}
	
	return &ast.BlockStatement{Statements: statements}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	return &ast.ExpressionStatement{Expression: expr}
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseInfixExpression()
}

func (p *Parser) parseInfixExpression() ast.Expression {
	left := p.parsePrimaryExpression()
	
	// Chercher un opérateur dans le token suivant
	if p.peekToken.Type == lexer.OPERATOR && 
	   (p.peekToken.Literal == "+" || p.peekToken.Literal == "-" || 
		p.peekToken.Literal == "*" || p.peekToken.Literal == "/" ||
		p.peekToken.Literal == "==" || p.peekToken.Literal == "!=" ||
		p.peekToken.Literal == "<" || p.peekToken.Literal == ">" ||
		p.peekToken.Literal == "<=" || p.peekToken.Literal == ">=" ||
		p.peekToken.Literal == "&&" || p.peekToken.Literal == "||") {
		
		p.nextToken() // aller sur l'opérateur
		operator := p.curToken.Literal
		p.nextToken() // aller sur l'opérande de droite
		right := p.parseInfixExpression()
		
		return &ast.InfixExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	
	return left
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.curToken.Type {
	case lexer.IDENT:
		// Gérer 'console' comme un identifiant spécial
		if p.curToken.Literal == "console" {
			return p.parseConsoleCall()
		}
		return p.parseIdentifierOrCall()
	case lexer.STRING:
		return &ast.StringLiteral{Value: p.curToken.Literal}
	case lexer.TEMPLATE:
		return p.parseTemplateLiteral()
	case lexer.NUMBER:
		return &ast.NumberLiteral{Value: p.curToken.Literal}
	case lexer.LBRACKET:
		return p.parseArrayLiteral()
	case lexer.LBRACE:
		return p.parseObjectLiteral()
	case lexer.KEYWORD:
		if p.curToken.Literal == "true" {
			return &ast.BooleanLiteral{Value: true}
		} else if p.curToken.Literal == "false" {
			return &ast.BooleanLiteral{Value: false}
		}
	}
	return nil
}

func (p *Parser) parseIdentifierOrCall() ast.Expression {
	ident := &ast.Identifier{Value: p.curToken.Literal}
	
	// Vérifier s'il y a un appel de fonction
	if p.peekToken.Type == lexer.LPAREN {
		p.nextToken() // aller à '('
		return p.parseFunctionCall(ident)
	}
	
	// Vérifier l'accès par index [0] ou propriété .prop
	if p.peekToken.Type == lexer.LBRACKET {
		p.nextToken() // aller à '['
		return p.parseIndexAccess(ident)
	}
	
	if p.peekToken.Type == lexer.DOT {
		p.nextToken() // aller à '.'
		return p.parseDotAccess(ident)
	}
	
	return ident
}

func (p *Parser) parseFunctionCall(fn ast.Expression) ast.Expression {
	p.nextToken() // passer '('
	
	var args []ast.Expression
	for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}
		
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else {
			break
		}
	}
	
	// Passer la parenthèse fermante
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken()
	}
	
	return &ast.CallExpression{Function: fn, Arguments: args}
}

func (p *Parser) parseIndexAccess(obj ast.Expression) ast.Expression {
	p.nextToken() // passer '['
	index := p.parseExpression()
	
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken() // passer ']'
	}
	
	return &ast.IndexExpression{Left: obj, Index: index}
}

func (p *Parser) parseDotAccess(obj ast.Expression) ast.Expression {
	p.nextToken() // passer '.'
	property := p.curToken.Literal
	p.nextToken()
	
	return &ast.DotExpression{Object: obj, Property: property}
}

func (p *Parser) parseConsoleCall() ast.Expression {
	// On est sur 'console'
	consoleObj := &ast.Identifier{Value: "console"}
	
	// Vérifier s'il y a un point
	if p.peekToken.Type == lexer.DOT {
		p.nextToken() // aller à '.'
		p.nextToken() // aller à 'log' (ou autre méthode)
		
		method := p.curToken.Literal
		methodAccess := &ast.DotExpression{Object: consoleObj, Property: method}
		
		// Vérifier s'il y a un appel de fonction
		if p.peekToken.Type == lexer.LPAREN {
			p.nextToken() // aller à '('
			return p.parseFunctionCall(methodAccess)
		}
		
		return methodAccess
	}
	
	return consoleObj
}

func (p *Parser) parseTemplateLiteral() ast.Expression {
	return &ast.TemplateLiteral{Parts: []ast.Expression{&ast.StringLiteral{Value: p.curToken.Literal}}}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	p.nextToken() // passer '['
	
	var elements []ast.Expression
	for p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
		element := p.parseExpression()
		if element != nil {
			elements = append(elements, element)
		}
		
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else if p.curToken.Type != lexer.RBRACKET {
			// Avancer si ce n'est ni une virgule ni la fin
			p.nextToken()
		}
	}
	
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken() // passer ']'
	}
	
	return &ast.ArrayLiteral{Elements: elements}
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	p.nextToken() // passer '{'
	
	var properties []ast.ObjectProperty
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENT {
			key := p.curToken.Literal
			p.nextToken() // aller à ':'
			
			if p.curToken.Type == lexer.COLON {
				p.nextToken() // passer ':'
				value := p.parseExpression()
				properties = append(properties, ast.ObjectProperty{Key: key, Value: value})
				
				if p.curToken.Type == lexer.COMMA {
					p.nextToken()
				}
			}
		} else {
			p.nextToken()
		}
	}
	
	if p.curToken.Type == lexer.RBRACE {
		p.nextToken() // passer '}'
	}
	
	return &ast.ObjectLiteral{Properties: properties}
}

func (p *Parser) skipUnsupportedStatement() {
	// Ignorer jusqu'au prochain ; ou } ou fin de fichier
	for p.curToken.Type != lexer.SEMICOLON && 
		p.curToken.Type != lexer.RBRACE && 
		p.curToken.Type != lexer.EOF {
		p.nextToken()
	}
}

func (p *Parser) parseTypeAlias() ast.Statement {
	// type TaskStatus = 'pending' | 'in_progress' | 'done';
	// Pour l'instant, on va juste capturer le nom et ignorer le reste
	p.nextToken() // passer 'type'
	name := p.curToken.Literal
	
	// Ignorer le reste jusqu'au ;
	for p.curToken.Type != lexer.SEMICOLON && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}
	
	return &ast.TypeAlias{Name: name, Type: "string"} // simplifié
}

func (p *Parser) parseInterface() ast.Statement {
	// interface Task { ... }
	p.nextToken() // passer 'interface'
	name := p.curToken.Literal
	
	// Ignorer le corps de l'interface pour l'instant
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}
	
	return &ast.Interface{Name: name, Fields: []ast.InterfaceField{}}
}

func (p *Parser) parseClass() ast.Statement {
	// class TaskManager { ... }
	p.nextToken() // passer 'class'
	name := p.curToken.Literal
	
	// Ignorer le corps de la classe pour l'instant
	braceCount := 0
	for {
		if p.curToken.Type == lexer.LBRACE {
			braceCount++
		} else if p.curToken.Type == lexer.RBRACE {
			braceCount--
			if braceCount == 0 {
				break
			}
		} else if p.curToken.Type == lexer.EOF {
			break
		}
		p.nextToken()
	}
	
	return &ast.ClassDeclaration{Name: name, Fields: []ast.ClassField{}, Methods: []ast.ClassMethod{}}
}

func (p *Parser) parseFunction() ast.Statement {
	// function name(params): returnType { body }
	p.nextToken() // passer 'function'
	
	if p.curToken.Type != lexer.IDENT {
		return nil
	}
	
	name := p.curToken.Literal
	p.nextToken() // aller à '('
	
	if p.curToken.Type != lexer.LPAREN {
		return nil
	}
	
	params := p.parseParameters()
	
	// Type de retour optionnel
	var returnType string
	if p.curToken.Type == lexer.COLON {
		p.nextToken() // passer ':'
		if p.curToken.Type == lexer.IDENT {
			returnType = p.curToken.Literal
			p.nextToken()
		}
	}
	
	// Corps de la fonction
	if p.curToken.Type != lexer.LBRACE {
		return nil
	}
	
	body := p.parseBlockStatement()
	
	var bodyStatements []ast.Statement
	if blockStmt, ok := body.(*ast.BlockStatement); ok {
		bodyStatements = blockStmt.Statements
	}
	
	return &ast.FunctionDeclaration{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       bodyStatements,
	}
}

func (p *Parser) parseParameters() []ast.Parameter {
	var params []ast.Parameter
	
	p.nextToken() // passer '('
	
	// Boucle avec sécurité contre les boucles infinies
	maxIterations := 50 // Sécurité
	iterations := 0
	
	for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF && iterations < maxIterations {
		iterations++
		
		if p.curToken.Type == lexer.IDENT {
			param := ast.Parameter{Name: p.curToken.Literal}
			p.nextToken()
			
			// Type optionnel
			if p.curToken.Type == lexer.COLON {
				p.nextToken() // passer ':'
				if p.curToken.Type == lexer.IDENT {
					param.Type = p.curToken.Literal
					p.nextToken()
				}
			}
			
			params = append(params, param)
			
			if p.curToken.Type == lexer.COMMA {
				p.nextToken() // passer ','
			} else if p.curToken.Type != lexer.RPAREN {
				// Si ce n'est ni ',' ni ')', sortir pour éviter boucle infinie
				break
			}
		} else {
			// Avancer si token inattendu pour éviter boucle infinie
			p.nextToken()
		}
	}
	
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // passer ')'
	}
	
	return params
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	vd := &ast.VariableDeclaration{}
	vd.IsConst = p.curToken.Literal == "const"

	p.nextToken() // identifiant
	vd.Name = p.curToken.Literal

	p.nextToken() // : ou =
	if p.curToken.Type == lexer.COLON {
		// Skip le type (simple ou complexe) jusqu'à =
		for p.curToken.Type != lexer.OPERATOR || p.curToken.Literal != "=" {
			if p.curToken.Type == lexer.EOF {
				break
			}
			
			// Capturer les types simples
			if p.curToken.Type == lexer.IDENT && vd.Type == "" {
				vd.Type = p.curToken.Literal
			}
			
			// Gérer les types array comme "number[]"
			if p.curToken.Type == lexer.LBRACKET && p.peekToken.Type == lexer.RBRACKET {
				vd.Type = vd.Type + "[]"
				p.nextToken() // skip ']'
			}
			
			p.nextToken()
		}
	}

	if p.curToken.Type == lexer.OPERATOR && p.curToken.Literal == "=" {
		p.nextToken()
		vd.Value = p.parseExpression()
	}

	return vd
}
func (p *Parser) ParseProgram() []ast.Statement {
	var statements []ast.Statement

	for p.curToken.Type != lexer.EOF {
		// Ignorer les commentaires
		if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
			continue
		}
		
		stmt := p.ParseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

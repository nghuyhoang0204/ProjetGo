package parser

import (
	"ProjetGo/ast"
	"ProjetGo/lexer"
	"fmt"
)

// Priorités des opérateurs
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > ou <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X ou !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

// Map des priorités des opérateurs
var precedences = map[lexer.TokenType]int{
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LTE:      LESSGREATER,
	lexer.GTE:      LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.DIVIDE:   PRODUCT,
	lexer.MULTIPLY: PRODUCT,
	lexer.MODULO:   PRODUCT,
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: INDEX,
	lexer.DOT:      INDEX,
	lexer.AND:      EQUALS,
	lexer.OR:       EQUALS,
}

// Fonctions de parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser analyseur syntaxique
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn

	errors []string
}

// New crée un nouveau parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialiser les fonctions de parsing
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TEMPLATE, p.parseTemplateLiteral)
	p.registerPrefix(lexer.KEYWORD, p.parseKeywordExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseObjectLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.DIVIDE, p.parseInfixExpression)
	p.registerInfix(lexer.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(lexer.MODULO, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseMemberExpression)

	// Lire deux tokens pour initialiser curToken et peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// registerPrefix enregistre une fonction de parsing préfixe
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix enregistre une fonction de parsing infixe
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken avance au token suivant
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	
	// Ignorer les commentaires
	for p.peekToken.Type == lexer.COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

// ParseProgram parse le programme complet
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.EOF {
		// Ignorer les commentaires
		if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parse un statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.KEYWORD:
		return p.parseKeywordStatement()
	case lexer.IDENT:
		// Vérifier si c'est une assignation
		if p.peekToken.Type == lexer.ASSIGN {
			return p.parseAssignmentStatement()
		}
		// Sinon c'est une expression statement
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseKeywordStatement parse les statements qui commencent par un mot-clé
func (p *Parser) parseKeywordStatement() ast.Statement {
	switch p.curToken.Literal {
	case "let", "const", "var":
		return p.parseVariableDeclaration()
	case "function":
		return p.parseFunctionDeclaration()
	case "if":
		return p.parseIfStatement()
	case "for":
		return p.parseForStatement()
	case "while":
		return p.parseWhileStatement()
	case "return":
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseVariableDeclaration parse: let/const/var name: type = value
func (p *Parser) parseVariableDeclaration() ast.Statement {
	stmt := &ast.VariableDeclaration{}
	stmt.Token = p.curToken.Literal
	stmt.IsConst = p.curToken.Literal == "const"

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	// Type optionnel
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // consommer ':'
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.Type = p.curToken.Literal
		
		// Gérer les types array comme number[]
		if p.peekToken.Type == lexer.LBRACKET {
			p.nextToken() // consommer '['
			if p.peekToken.Type == lexer.RBRACKET {
				p.nextToken() // consommer ']'
				stmt.Type += "[]"
			}
		}
	}

	// Valeur optionnelle
	if p.peekToken.Type == lexer.ASSIGN {
		p.nextToken() // consommer '='
		p.nextToken() // aller à la valeur
		stmt.Value = p.parseExpression(LOWEST)
	}

	// Consommer le point-virgule optionnel
	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseFunctionDeclaration parse: function name(params): returnType { body }
func (p *Parser) parseFunctionDeclaration() ast.Statement {
	stmt := &ast.FunctionDeclaration{}
	stmt.Token = p.curToken.Literal

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Type de retour optionnel
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // consommer ':'
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.ReturnType = p.curToken.Literal
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parse les paramètres d'une fonction
func (p *Parser) parseFunctionParameters() []ast.Parameter {
	identifiers := []ast.Parameter{}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	param := ast.Parameter{}
	param.Name = p.curToken.Literal

	// Type optionnel
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // consommer ':'
		if !p.expectPeek(lexer.IDENT) {
			return identifiers
		}
		param.Type = p.curToken.Literal
	}

	identifiers = append(identifiers, param)

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()

		param := ast.Parameter{}
		param.Name = p.curToken.Literal

		// Type optionnel
		if p.peekToken.Type == lexer.COLON {
			p.nextToken() // consommer ':'
			if !p.expectPeek(lexer.IDENT) {
				return identifiers
			}
			param.Type = p.curToken.Literal
		}

		identifiers = append(identifiers, param)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return identifiers
}

// parseIfStatement parse: if (condition) { then } else { else }
func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{}
	stmt.Token = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.ThenBranch = p.parseBlockStatement()

	if p.peekToken.Type == lexer.KEYWORD && p.peekToken.Literal == "else" {
		p.nextToken()

		if p.peekToken.Type == lexer.LBRACE {
			p.nextToken()
			stmt.ElseBranch = p.parseBlockStatement()
		}
	}

	return stmt
}

// parseForStatement parse: for (init; condition; update) { body }
func (p *Parser) parseForStatement() ast.Statement {
	stmt := &ast.ForStatement{}
	stmt.Token = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Init
	p.nextToken()
	if p.curToken.Type != lexer.SEMICOLON {
		stmt.Init = p.parseStatement()
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Condition
	p.nextToken()
	if p.curToken.Type != lexer.SEMICOLON {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	// Update
	p.nextToken()
	if p.curToken.Type != lexer.RPAREN {
		stmt.Update = p.parseExpressionStatement()
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseWhileStatement parse: while (condition) { body }
func (p *Parser) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStatement{}
	stmt.Token = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseReturnStatement parse: return expression;
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{}
	stmt.Token = p.curToken.Literal

	p.nextToken()

	if p.curToken.Type == lexer.SEMICOLON {
		return stmt
	}

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parse: { statements... }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Token = p.curToken.Literal
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		// Ignorer les commentaires
		if p.curToken.Type == lexer.COMMENT {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseAssignmentStatement parse: variable = value
func (p *Parser) parseAssignmentStatement() ast.Statement {
	stmt := &ast.AssignmentStatement{}
	stmt.Token = p.curToken.Literal
	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parse une expression utilisée comme statement
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{}
	stmt.Token = p.curToken.Literal
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseExpression parse une expression avec la priorité donnée
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != lexer.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Parse functions for expressions

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken.Literal, Value: p.curToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{}
	lit.Token = p.curToken.Literal
	lit.Value = p.curToken.Literal
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken.Literal, Value: p.curToken.Literal}
}

func (p *Parser) parseTemplateLiteral() ast.Expression {
	return &ast.TemplateLiteral{
		Token: p.curToken.Literal,
		Parts: []ast.Expression{&ast.StringLiteral{Value: p.curToken.Literal}},
	}
}

func (p *Parser) parseKeywordExpression() ast.Expression {
	switch p.curToken.Literal {
	case "true":
		return &ast.BooleanLiteral{Token: p.curToken.Literal, Value: true}
	case "false":
		return &ast.BooleanLiteral{Token: p.curToken.Literal, Value: false}
	default:
		return &ast.Identifier{Token: p.curToken.Literal, Value: p.curToken.Literal}
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.InfixExpression{}
	expression.Token = p.curToken.Literal
	expression.Operator = p.curToken.Literal

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{}
	expression.Token = p.curToken.Literal
	expression.Left = left
	expression.Operator = p.curToken.Literal

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken.Literal}
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken.Literal}
	obj.Properties = []ast.ObjectProperty{}

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		return obj
	}

	p.nextToken()

	for {
		prop := ast.ObjectProperty{}

		if p.curToken.Type == lexer.IDENT {
			prop.Key = p.curToken.Literal
		} else if p.curToken.Type == lexer.STRING {
			prop.Key = p.curToken.Literal
		} else {
			return nil
		}

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		prop.Value = p.parseExpression(LOWEST)

		obj.Properties = append(obj.Properties, prop)

		if p.peekToken.Type != lexer.COMMA {
			break
		}

		p.nextToken()
		p.nextToken()
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return obj
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken.Literal, Function: fn}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{}
	exp.Token = p.curToken.Literal
	exp.Object = left
	exp.Computed = true

	p.nextToken()
	exp.Property = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{}
	exp.Token = p.curToken.Literal
	exp.Object = left
	exp.Computed = false

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	exp.Property = &ast.Identifier{Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// Helper functions

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Error handling

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

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
	if p.curToken.Literal == "let" || p.curToken.Literal == "const" {
		return p.parseVariableDeclaration()
	}
	// à compléter : function, return, etc.
	return nil
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	vd := &ast.VariableDeclaration{}
	vd.IsConst = p.curToken.Literal == "const"

	p.nextToken() // identifiant
	vd.Name = p.curToken.Literal

	p.nextToken() // :
	if p.curToken.Type == lexer.COLON {
		p.nextToken() // type (ex: "string")
		vd.Type = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == lexer.OPERATOR && p.curToken.Literal == "=" {
		p.nextToken()
		switch p.curToken.Type {
		case lexer.STRING:
			vd.Value = &ast.StringLiteral{Value: p.curToken.Literal}
		case lexer.NUMBER:
			vd.Value = &ast.NumberLiteral{Value: p.curToken.Literal}
		}
	}

	return vd
}
func (p *Parser) ParseProgram() []ast.Statement {
	var statements []ast.Statement

	for p.curToken.Type != lexer.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

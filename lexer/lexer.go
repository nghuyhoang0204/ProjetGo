package lexer

import (
	"unicode"
)

type TokenType int

const (
	IDENTIFIER TokenType = iota
	KEYWORD
	OPERATOR
	NUMBER
	STRING
	EOF
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Tokenize converts input TypeScript code into tokens.
func Tokenize(input string) []Token {
	var tokens []Token
	var currentToken string
	var currentLine, currentColumn int

	for _, r := range input {
		currentColumn++
		if unicode.IsSpace(r) {
			if currentToken != "" {
				tokens = append(tokens, Token{Type: IDENTIFIER, Value: currentToken, Line: currentLine, Column: currentColumn})
				currentToken = ""
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			currentToken += string(r)
		} else {
			if currentToken != "" {
				tokens = append(tokens, Token{Type: IDENTIFIER, Value: currentToken, Line: currentLine, Column: currentColumn})
				currentToken = ""
			}
			// Handle special characters like ':' and ';'
			tokens = append(tokens, Token{Type: OPERATOR, Value: string(r), Line: currentLine, Column: currentColumn})
		}
	}

	if currentToken != "" {
		tokens = append(tokens, Token{Type: IDENTIFIER, Value: currentToken, Line: currentLine, Column: currentColumn})
	}

	return tokens
}

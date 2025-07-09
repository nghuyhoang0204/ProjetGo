package lexer

// 🚀 Analyseur lexical TypeScript → JavaScript (Go pur)
// Convertit le code source en tokens pour le parser

// TokenType représente les différents types de tokens
type TokenType string

// Token représente un token avec son type et sa valeur
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Types de tokens
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiants et littéraux
	IDENT    = "IDENT"    // variables, fonctions
	NUMBER   = "NUMBER"   // 123, 45.67
	STRING   = "STRING"   // "hello", 'world'
	TEMPLATE = "TEMPLATE" // `template string`
	COMMENT  = "COMMENT"  // // ou /* */

	// Mots-clés
	KEYWORD = "KEYWORD"

	// Opérateurs
	ASSIGN   = "="     // =
	PLUS     = "+"     // +
	MINUS    = "-"     // -
	MULTIPLY = "*"     // *
	DIVIDE   = "/"     // /
	MODULO   = "%"     // %
	
	// Opérateurs de comparaison
	EQ     = "=="    // ==
	NOT_EQ = "!="    // !=
	LT     = "<"     // <
	GT     = ">"     // >
	LTE    = "<="    // <=
	GTE    = ">="    // >=
	
	// Opérateurs logiques
	AND = "&&"   // &&
	OR  = "||"   // ||
	NOT = "!"    // !
	
	// Opérateurs d'assignation
	PLUS_ASSIGN     = "+="   // +=
	MINUS_ASSIGN    = "-="   // -=
	MULTIPLY_ASSIGN = "*="   // *=
	DIVIDE_ASSIGN   = "/="   // /=
	
	// Incrémentation/décrémentation
	INCREMENT = "++"   // ++
	DECREMENT = "--"   // --

	// Délimiteurs
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."
	QUESTION  = "?"

	// Parenthèses et crochets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Symboles TypeScript
	ARROW      = "=>"   // =>
	OPTIONAL   = "?:"   // ?:
	PIPE       = "|"    // |
	AMPERSAND  = "&"    // &
	ELLIPSIS   = "..."  // ...
)

// Mots-clés TypeScript/JavaScript
var keywords = map[string]TokenType{
	// Déclarations
	"let":       KEYWORD,
	"const":     KEYWORD,
	"var":       KEYWORD,
	"function":  KEYWORD,
	"class":     KEYWORD,
	"interface": KEYWORD,
	"type":      KEYWORD,
	"enum":      KEYWORD,
	"namespace": KEYWORD,

	// Contrôle de flux
	"if":       KEYWORD,
	"else":     KEYWORD,
	"for":      KEYWORD,
	"while":    KEYWORD,
	"do":       KEYWORD,
	"switch":   KEYWORD,
	"case":     KEYWORD,
	"default":  KEYWORD,
	"break":    KEYWORD,
	"continue": KEYWORD,
	"return":   KEYWORD,

	// Valeurs littérales
	"true":      KEYWORD,
	"false":     KEYWORD,
	"null":      KEYWORD,
	"undefined": KEYWORD,

	// Fonctions/objets
	"new":    KEYWORD,
	"this":   KEYWORD,
	"super":  KEYWORD,
	"static": KEYWORD,

	// Modificateurs
	"public":    KEYWORD,
	"private":   KEYWORD,
	"protected": KEYWORD,
	"readonly":  KEYWORD,
	"abstract":  KEYWORD,

	// Async/await
	"async": KEYWORD,
	"await": KEYWORD,

	// Modules
	"import":  KEYWORD,
	"export":  KEYWORD,
	"from":    KEYWORD,
	"as":      KEYWORD,

	// Gestion d'erreurs
	"try":     KEYWORD,
	"catch":   KEYWORD,
	"throw":   KEYWORD,
	"finally": KEYWORD,

	// Types
	"string":  KEYWORD,
	"number":  KEYWORD,
	"boolean": KEYWORD,
	"object":  KEYWORD,
	"any":     KEYWORD,
	"void":    KEYWORD,
	"never":   KEYWORD,

	// Autres
	"typeof":     KEYWORD,
	"instanceof": KEYWORD,
	"in":         KEYWORD,
	"of":         KEYWORD,
	"extends":    KEYWORD,
	"implements": KEYWORD,
}

// Lexer analyseur lexical
type Lexer struct {
	input        string
	position     int  // position actuelle dans input
	readPosition int  // position de lecture suivante
	ch           byte // caractère courant sous examen
	line         int  // ligne actuelle
	column       int  // colonne actuelle
}

// New crée un nouveau lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar lit le caractère suivant et avance dans l'input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL représente "EOF"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar retourne le caractère suivant sans avancer
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekCharAt retourne le caractère à la position relative
func (l *Lexer) peekCharAt(offset int) byte {
	pos := l.readPosition + offset - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// NextToken génère le token suivant
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: ARROW, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: INCREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PLUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(PLUS, l.ch, l.line, l.column)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DECREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MINUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(MINUS, l.ch, l.line, l.column)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MULTIPLY_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(MULTIPLY, l.ch, l.line, l.column)
		}
	case '/':
		if l.peekChar() == '/' {
			tok.Type = COMMENT
			tok.Literal = l.readLineComment()
			tok.Line = l.line
			tok.Column = l.column
		} else if l.peekChar() == '*' {
			tok.Type = COMMENT
			tok.Literal = l.readBlockComment()
			tok.Line = l.line
			tok.Column = l.column
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DIVIDE_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(DIVIDE, l.ch, l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(NOT, l.ch, l.line, l.column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(GT, l.ch, l.line, l.column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(AMPERSAND, l.ch, l.line, l.column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(PIPE, l.ch, l.line, l.column)
		}
	case '?':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OPTIONAL, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(QUESTION, l.ch, l.line, l.column)
		}
	case '.':
		if l.peekChar() == '.' && l.peekCharAt(2) == '.' {
			l.readChar()
			l.readChar()
			tok = Token{Type: ELLIPSIS, Literal: "...", Line: l.line, Column: l.column}
		} else {
			tok = newToken(DOT, l.ch, l.line, l.column)
		}
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column)
	case ':':
		tok = newToken(COLON, l.ch, l.line, l.column)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, l.column)
	case '%':
		tok = newToken(MODULO, l.ch, l.line, l.column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
		tok.Line = l.line
		tok.Column = l.column
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
		tok.Line = l.line
		tok.Column = l.column
	case '`':
		tok.Type = TEMPLATE
		tok.Literal = l.readTemplateString()
		tok.Line = l.line
		tok.Column = l.column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

// newToken crée un nouveau token
func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

// readIdentifier lit un identifiant
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber lit un nombre (entier ou décimal)
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// Nombre décimal
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString lit une chaîne de caractères
func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		// Gérer les caractères d'échappement
		if l.ch == '\\' {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// readTemplateString lit une template string
func (l *Lexer) readTemplateString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '`' || l.ch == 0 {
			break
		}
		// Pour simplifier, on ne gère pas les expressions ${} pour l'instant
		if l.ch == '\\' {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// readLineComment lit un commentaire de ligne //
func (l *Lexer) readLineComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readBlockComment lit un commentaire de bloc /* */
func (l *Lexer) readBlockComment() string {
	position := l.position
	for {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			break
		}
		if l.ch == 0 {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhitespace ignore les espaces, tabs, retours à la ligne
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter vérifie si c'est une lettre
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// isDigit vérifie si c'est un chiffre
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent détermine si c'est un mot-clé ou un identifiant
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// IsKeyword vérifie si une chaîne est un mot-clé
func IsKeyword(literal string) bool {
	_, ok := keywords[literal]
	return ok
}

// GetKeywordType retourne le type de token pour un mot-clé
func GetKeywordType(literal string) TokenType {
	if tok, ok := keywords[literal]; ok {
		return tok
	}
	return IDENT
}

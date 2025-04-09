package lexer

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"

    IDENT     = "IDENT"     // variable, fonction
    KEYWORD   = "KEYWORD"   // let, const, function, return
    NUMBER    = "NUMBER"    // 123, 3.14
    STRING    = "STRING"    // "hello"
    OPERATOR  = "OPERATOR"  // =, +, -, *, /
    COLON     = ":"
    SEMICOLON = ";"
    COMMA     = ","
    LPAREN    = "("
    RPAREN    = ")"
    LBRACE    = "{"
    RBRACE    = "}"
)

var keywords = map[string]TokenType{
    "let":      KEYWORD,
    "const":    KEYWORD,
    "var":      KEYWORD,
    "function": KEYWORD,
    "return":   KEYWORD,
}

type Lexer struct {
    input        string
    position     int  // position actuelle
    readPosition int  // position après lecture
    ch           byte // caractère courant
    line         int
    column       int
}

func New(input string) *Lexer {
    l := &Lexer{input: input, line: 1, column: 0}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    if l.ch == '\n' {
        l.line++
        l.column = 0
    } else {
        l.column++
    }
    l.position = l.readPosition
    l.readPosition++
}

func (l *Lexer) NextToken() Token {
    l.skipWhitespace()

    tok := Token{Line: l.line, Column: l.column}

    switch l.ch {
    case '=':
        tok = newToken(OPERATOR, "=", l)
    case '+':
        tok = newToken(OPERATOR, "+", l)
    case ';':
        tok = newToken(SEMICOLON, ";", l)
    case ':':
        tok = newToken(COLON, ":", l)
    case '(':
        tok = newToken(LPAREN, "(", l)
    case ')':
        tok = newToken(RPAREN, ")", l)
    case '{':
        tok = newToken(LBRACE, "{", l)
    case '}':
        tok = newToken(RBRACE, "}", l)
    case '"':
        tok.Type = STRING
        tok.Literal = l.readString()
    case 0:
        tok.Type = EOF
        tok.Literal = ""
    default:
        if isLetter(l.ch) {
            ident := l.readIdentifier()
            tok.Literal = ident
            if t, ok := keywords[ident]; ok {
                tok.Type = t
            } else {
                tok.Type = IDENT
            }
            return tok
        } else if isDigit(l.ch) {
            tok.Type = NUMBER
            tok.Literal = l.readNumber()
            return tok
        } else {
            tok = newToken(ILLEGAL, string(l.ch), l)
        }
    }

    l.readChar()
    return tok
}

func newToken(tokenType TokenType, ch string, l *Lexer) Token {
    return Token{Type: tokenType, Literal: ch, Line: l.line, Column: l.column}
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readIdentifier() string {
    start := l.position
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
    start := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[start:l.position]
}

func (l *Lexer) readString() string {
    l.readChar() // skip "
    start := l.position
    for l.ch != '"' && l.ch != 0 {
        l.readChar()
    }
    str := l.input[start:l.position]
    l.readChar() // skip closing "
    return str
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

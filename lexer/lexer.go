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
    TEMPLATE  = "TEMPLATE"  // `hello ${name}`
    COMMENT   = "COMMENT"   // // commentaire
    OPERATOR  = "OPERATOR"  // =, +, -, *, /
    COLON     = ":"
    SEMICOLON = ";"
    COMMA     = ","
    LPAREN    = "("
    RPAREN    = ")"
    LBRACE    = "{"
    RBRACE    = "}"
    LBRACKET  = "["
    RBRACKET  = "]"
    DOT       = "."
    PIPE      = "|"
    ARROW     = "=>"
    QUESTION  = "?"
    EXCLAMATION = "!"
)

var keywords = map[string]TokenType{
    "let":       KEYWORD,
    "const":     KEYWORD,
    "var":       KEYWORD,
    "function":  KEYWORD,
    "return":    KEYWORD,
    "true":      KEYWORD,
    "false":     KEYWORD,
    "type":      KEYWORD,
    "interface": KEYWORD,
    "class":     KEYWORD,
    "async":     KEYWORD,
    "await":     KEYWORD,
    "new":       KEYWORD,
    "this":      KEYWORD,
    "private":   KEYWORD,
    "public":    KEYWORD,
    "static":    KEYWORD,
    "if":        KEYWORD,
    "else":      KEYWORD,
    "for":       KEYWORD,
    "while":     KEYWORD,
    "do":        KEYWORD,
    "break":     KEYWORD,
    "continue":  KEYWORD,
    "try":       KEYWORD,
    "catch":     KEYWORD,
    "throw":     KEYWORD,
    "switch":    KEYWORD,
    "case":      KEYWORD,
    "default":   KEYWORD,
    "console":   KEYWORD,
    "void":      KEYWORD,
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

func (l *Lexer) peekChar() byte {
    if l.readPosition >= len(l.input) {
        return 0
    } else {
        return l.input[l.readPosition]
    }
}

func (l *Lexer) NextToken() Token {
    l.skipWhitespace()

    tok := Token{Line: l.line, Column: l.column}

    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else if l.peekChar() == '>' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: ARROW, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, "=", l)
        }
    case '+':
        if l.peekChar() == '+' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, "+", l)
        }
    case '-':
        if l.peekChar() == '-' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, "-", l)
        }
    case '*':
        tok = newToken(OPERATOR, "*", l)
    case '/':
        if l.peekChar() == '/' {
            tok.Type = COMMENT
            tok.Literal = l.readComment()
        } else {
            tok = newToken(OPERATOR, "/", l)
        }
    case '<':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, "<", l)
        }
    case '>':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, ">", l)
        }
    case '!':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(EXCLAMATION, "!", l)
        }
    case '&':
        if l.peekChar() == '&' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(OPERATOR, "&", l)
        }
    case '|':
        if l.peekChar() == '|' {
            ch := l.ch
            l.readChar()
            tok = Token{Type: OPERATOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
        } else {
            tok = newToken(PIPE, "|", l)
        }
    case '?':
        tok = newToken(QUESTION, "?", l)
    case '.':
        tok = newToken(DOT, ".", l)
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
    case '[':
        tok = newToken(LBRACKET, "[", l)
    case ']':
        tok = newToken(RBRACKET, "]", l)
    case ',':
        tok = newToken(COMMA, ",", l)
    case '"':
        tok.Type = STRING
        tok.Literal = l.readString()
    case '\'':
        tok.Type = STRING
        tok.Literal = l.readSingleQuoteString()
    case '`':
        tok.Type = TEMPLATE
        tok.Literal = l.readTemplateLiteral()
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

func (l *Lexer) readSingleQuoteString() string {
    l.readChar() // skip '
    start := l.position
    for l.ch != '\'' && l.ch != 0 {
        l.readChar()
    }
    str := l.input[start:l.position]
    l.readChar() // skip closing '
    return str
}

func (l *Lexer) readTemplateLiteral() string {
    l.readChar() // skip `
    start := l.position
    for l.ch != '`' && l.ch != 0 {
        l.readChar()
    }
    str := l.input[start:l.position]
    l.readChar() // skip closing `
    return str
}

func (l *Lexer) readComment() string {
    start := l.position
    for l.ch != '\n' && l.ch != 0 {
        l.readChar()
    }
    return l.input[start:l.position]
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

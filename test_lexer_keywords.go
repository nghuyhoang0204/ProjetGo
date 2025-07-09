package main

import (
	"ProjetGo/lexer"
	"fmt"
)

func testLexerKeywords() {
	code := `if (age >= 18) {
  majeur = true;
} else {
  majeur = false;
}

for (let i = 0; i < notes.length; i++) {
  console.log("Note");
}

while (compteur > 0) {
  compteur--;
}`

	fmt.Println("=== TEST LEXER KEYWORDS ===")
	fmt.Println("Code:", code)
	fmt.Println()

	l := lexer.New(code)
	
	for i := 0; i < 50; i++ { // Limiter pour éviter boucle infinie
		tok := l.NextToken()
		if tok.Type == lexer.EOF {
			break
		}
		fmt.Printf("Token %d: Type=%s, Literal='%s'\n", i+1, tok.Type, tok.Literal)
	}
}

func main() {
	testLexerKeywords()
}

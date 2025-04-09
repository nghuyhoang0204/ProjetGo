package main

import (
	"ProjetGo/lexer"
    "ProjetGo/parser"
    "ProjetGo/generator"
    "fmt"
)

func main() {
	input := `
	const message: string = "Hello";
	let count: number = 42;
	`

	l := lexer.New(input)         
	p := parser.New(l)            
	program := p.ParseProgram() 

	output := generator.GenerateJS(program)

	fmt.Println("=== JavaScript Output ===")
	fmt.Println(output)
}

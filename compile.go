package main

import "fmt"

func compile(source string) {
	initScanner(source)
	line := -1
	scanner.Source = source
	for {
		token := scanToken()
		if token.Line != line {
			fmt.Printf("%4d ", token.Line)
			line = token.Line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.Type, source[token.Start:token.Start+token.Length])

		if token.Type == TOKEN_EOF {
			break
		}
	}
}

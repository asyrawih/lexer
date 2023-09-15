package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	l := NewLexer(r)

	expr := parseExpression(l)
	fmt.Printf("expr will evaluate: %v\n", expr)
	result, err := evaluateExpression(expr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("result: %v\n", result)
}

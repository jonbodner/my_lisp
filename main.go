package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	done := false
	depth := 0
	tokens := make([]Token, 0)
	for !done {
		line, _, err := bio.ReadLine()
		if err != nil {
			fmt.Errorf("error: %v", err)
			return
		}
		newTokens, newDepth := Scan(string(line))
		depth = depth + newDepth
		if depth < 0 {
			fmt.Println("Invalid -- Too many closing parens")
			depth = 0
			tokens = make([]Token, 0)
			continue
		}
		tokens = append(tokens, newTokens...)
		if depth == 0 {
			fmt.Println(tokens)
			expr, pos, err := Parse(tokens)
			fmt.Println(expr)
			fmt.Println(pos)
			fmt.Println(err)
			result, err := Eval(expr)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
			tokens = make([]Token, 0)
		}
	}
}

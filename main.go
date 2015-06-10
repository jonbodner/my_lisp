package main

import (
	"bufio"
	"fmt"
	"github.com/jonbodner/my_lisp/evaluator"
	"github.com/jonbodner/my_lisp/parser"
	"github.com/jonbodner/my_lisp/scanner"
	"github.com/jonbodner/my_lisp/types"
	"os"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	done := false
	depth := 0
	tokens := []types.Token{}
	for !done {
		line, _, err := bio.ReadLine()
		if err != nil {
			fmt.Errorf("error: %v", err)
			return
		}
		newTokens, newDepth := scanner.Scan(string(line))
		depth = depth + newDepth
		if depth < 0 {
			fmt.Println("Invalid -- Too many closing parens")
			depth = 0
			tokens = make([]types.Token, 0)
			continue
		}
		tokens = append(tokens, newTokens...)
		if depth == 0 {
			//fmt.Println(tokens)
			expr, _, err := parser.Parse(tokens)
			//fmt.Println(expr)
			//fmt.Println(pos)
			//fmt.Println(err)
			if err != nil {
				fmt.Println(err)
			} else {
				result, err := evaluator.Eval(expr)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(result)
				}
			}
			tokens = make([]types.Token, 0)
		}
	}
}

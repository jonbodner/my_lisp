package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jonbodner/my_lisp/evaluator"
	"github.com/jonbodner/my_lisp/global"
	"github.com/jonbodner/my_lisp/parser"
	"github.com/jonbodner/my_lisp/scanner"
	"github.com/jonbodner/my_lisp/types"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	done := false
	depth := 0
	tokens := []types.Token{}
	for !done {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			done = true
			continue
		}
		newTokens, newDepth := scanner.Scan(string(line))
		depth = depth + newDepth
		if depth < 0 {
			fmt.Println("Invalid -- Too many closing parens")
			depth = 0
			tokens = []types.Token{}
			continue
		}
		tokens = append(tokens, newTokens...)
		if depth == 0 {
			//global.Log(tokens)
			expr, _, err := parser.Parse(tokens)
			//global.Log(expr)
			//global.Log(pos)
			//global.Log(err)
			if err != nil {
				global.Log(err)
			} else {
				result, err := evaluator.Eval(expr)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(result)
				}
			}
			tokens = []types.Token{}
		}
	}
}

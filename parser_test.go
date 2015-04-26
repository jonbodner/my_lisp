package main

import (
	"fmt"
	"testing"
)

func TestParserEmpty(t *testing.T) {
	in := ""
	fmt.Println(getExpression(in))
}

func TestParserAtom(t *testing.T) {
	in := "hello"
	fmt.Println(getExpression(in))
}

func TestParserEmptyList(t *testing.T) {
	in := "()"
	fmt.Println(getExpression(in))
}

func TestParserSimpleList(t *testing.T) {
	fmt.Println(getExpression("(a)"))
	fmt.Println(getExpression("(a b)"))
	fmt.Println(getExpression("(a b c)"))
}

func getExpression(in string) (Expr, int, error) {
	tokens, _ := Scan(in)
	expression, pos, err := Parse(tokens)
	return expression, pos, err
}


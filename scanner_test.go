package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestScannerEmpty(t *testing.T) {
	empty := ""

	tokens, depth := Scan(empty)

	testingHelper(t, []reflect.Type{}, 0, tokens, depth)
}

func TestScannerAtom(t *testing.T) {
	atom := "WORD"

	tokens, depth := Scan(atom)

	testingHelper(t, []reflect.Type{reflect.TypeOf(NAME(""))}, 0, tokens, depth)
}

func TestScannerList(t *testing.T) {
	atom := "(WORD WORD2 WORD3)"

	tokens, depth := Scan(atom)

	testingHelper(t,
		[]reflect.Type{
			reflect.TypeOf(LPAREN),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(RPAREN)},
		0, tokens, depth)
}

func TestScannerAllTokens(t *testing.T) {
	atom := "(WORD WORD2 WORD3 . '(1 2 3)) ("

	tokens, depth := Scan(atom)

	testingHelper(t,
		[]reflect.Type{
			reflect.TypeOf(LPAREN),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(DOT),
			reflect.TypeOf(QUOTE),
			reflect.TypeOf(LPAREN),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(NAME("")),
			reflect.TypeOf(RPAREN),
			reflect.TypeOf(RPAREN),
			reflect.TypeOf(LPAREN)},
		1, tokens, depth)
}

func testingHelper(t *testing.T, expectedTokens []reflect.Type, expectedDepth int, tokens []Token, depth int) {
	fmt.Println(tokens, depth)

	if len(tokens) != len(expectedTokens) {
		t.Errorf("Should have %d token(s), had %d", len(expectedTokens), len(tokens))
	}
	if depth != expectedDepth {
		t.Errorf("Should have depth of %d, had %d", expectedDepth, depth)
	}

	for pos, tokenType := range expectedTokens {
		if reflect.TypeOf(tokens[pos]) != tokenType {
			t.Errorf("Should be a %s, was a %s", tokenType.Name() , reflect.TypeOf(tokens[pos]).Name())
		}
	}
}

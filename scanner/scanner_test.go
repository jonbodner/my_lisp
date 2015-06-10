package scanner

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/jonbodner/my_lisp/types"
)

func TestScannerEmpty(t *testing.T) {
	empty := ""

	tokens, depth := Scan(empty)

	testingHelper(t, []reflect.Type{}, 0, tokens, depth)
}

func TestScannerAtom(t *testing.T) {
	atom := "WORD"

	tokens, depth := Scan(atom)

	testingHelper(t, []reflect.Type{reflect.TypeOf(types.NAME(""))}, 0, tokens, depth)
}

func TestScannerList(t *testing.T) {
	atom := "(WORD WORD2 WORD3)"

	tokens, depth := Scan(atom)

	testingHelper(t,
		[]reflect.Type{
			reflect.TypeOf(types.LPAREN),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.RPAREN)},
		0, tokens, depth)
}

func TestScannerAllTokens(t *testing.T) {
	atom := "(WORD WORD2 WORD3 . '(1 2 3)) ("

	tokens, depth := Scan(atom)

	testingHelper(t,
		[]reflect.Type{
			reflect.TypeOf(types.LPAREN),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.DOT),
			reflect.TypeOf(types.QUOTE),
			reflect.TypeOf(types.LPAREN),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.NAME("")),
			reflect.TypeOf(types.RPAREN),
			reflect.TypeOf(types.RPAREN),
			reflect.TypeOf(types.LPAREN)},
		1, tokens, depth)
}

func testingHelper(t *testing.T, expectedTokens []reflect.Type, expectedDepth int, tokens []types.Token, depth int) {
	fmt.Println(tokens, depth)

	if len(tokens) != len(expectedTokens) {
		t.Errorf("Should have %d token(s), had %d", len(expectedTokens), len(tokens))
	}
	if depth != expectedDepth {
		t.Errorf("Should have depth of %d, had %d", expectedDepth, depth)
	}

	for pos, tokenType := range expectedTokens {
		if reflect.TypeOf(tokens[pos]) != tokenType {
			t.Errorf("Should be a %s, was a %s", tokenType.Name(), reflect.TypeOf(tokens[pos]).Name())
		}
	}
}

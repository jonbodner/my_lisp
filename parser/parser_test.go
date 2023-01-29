package parser

import (
	"fmt"
	"github.com/jonbodner/my_lisp/assert"
	"github.com/jonbodner/my_lisp/scanner"
	"github.com/jonbodner/my_lisp/types"
	"testing"
)

func TestParserEmpty(t *testing.T) {
	a := assert.Assert{T: t}
	_, _, err := getExpression("")
	a.NotNil("err should have a value", err)
	a.Equals("wrong error message", "No tokens supplied", err.Error())
}

func TestParserAtom(t *testing.T) {
	a := assert.Assert{T: t}
	expr, _, err := getExpression("hello")
	a.Nil("err should not have a value", err)
	a.Equals("expected an atom", types.Atom("hello"), expr)
}

func TestParserEmptyList(t *testing.T) {
	a := assert.Assert{T: t}
	expr, _, err := getExpression("()")
	a.Nil("err should not have a value", err)
	s, ok := expr.(*types.SExpr)
	a.True("should be an SExpr", ok)
	l, _ := s.Left.(types.Nil)
	a.Equals("should be Nil", types.NIL, l)
	r, _ := s.Right.(types.Nil)
	a.Equals("should be Nil", types.NIL, r)
}

func TestBadLeftDottedPair(t *testing.T) {
	a := assert.Assert{T: t}
	_, _, err := getExpression("( . b)")
	a.NotNil("err should have a value", err)
	a.Equals("wrong error message", "Dot in unexpected location: ( _._ b ) ", err.Error())
}

func TestBadRightDottedPair(t *testing.T) {
	a := assert.Assert{T: t}
	_, _, err := getExpression("(a . )")
	a.NotNil("err should have a value", err)
	a.Equals("wrong error message", "Right paren in unexpected location: ( a . _)_ ", err.Error())
}

func TestBadEmptyDottedPair(t *testing.T) {
	a := assert.Assert{T: t}
	_, _, err := getExpression("( . )")
	a.NotNil("err should have a value", err)
	a.Equals("wrong error message", "Dot in unexpected location: ( _._ ) ", err.Error())
}

func TestGoodSimpleDottedPair(t *testing.T) {
	a := assert.Assert{T: t}
	expr, _, err := getExpression("( a . b)")
	a.Nil("err should not have a value", err)
	s, ok := expr.(*types.SExpr)
	a.True("should be an SExpr", ok)
	l, _ := s.Left.(types.Atom)
	a.Equals("should be an atom == a", types.Atom("a"), l)
	r, _ := s.Right.(types.Atom)
	a.Equals("should be an atom == b", types.Atom("b"), r)
}

func TestParserSimpleList(t *testing.T) {
	fmt.Println(getExpression("(a)"))
	fmt.Println(getExpression("(a b)"))
	fmt.Println(getExpression("(a b c)"))
}

func TestNestedList(t *testing.T) {
	a := assert.Assert{T: t}
	expr, pos, err := getExpression("((a b (c) d) (e f) g)")
	fmt.Println(expr, pos, err)
	a.Nil("err should not have a value", err)
}

func TestQuote(t *testing.T) {
	a := assert.Assert{T: t}
	expr, pos, err := getExpression("'a")
	fmt.Println(expr, pos, err)
	a.Nil("err should not have a value", err)

	expr, pos, err = getExpression("'(a b c)")
	fmt.Println(expr, pos, err)
	a.Nil("err should not have a value", err)
}

func TestQuoteNested(t *testing.T) {
	a := assert.Assert{T: t}
	expr, pos, err := getExpression("('(a b '(c) d) (e 'f) g)")
	fmt.Println(expr, pos, err)
	a.Nil("err should not have a value", err)
}

func getExpression(in string) (types.Expr, int, error) {
	tokens, _ := scanner.Scan(in)
	expression, pos, err := Parse(tokens)
	return expression, pos, err
}

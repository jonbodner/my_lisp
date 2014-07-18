package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"unicode"
)

//expressions
type SExpr struct {
	Left  Expr
	Right Expr
}

type Expr interface {
	isExpr()
}

func (s SExpr) isExpr() {}
func (s SExpr) String() string {
	return s.stringInner(true)
}

func (s SExpr) stringInner(top bool) string {
	out := ""
	if top {
		out += "("
	}
	switch l := s.Left.(type) {
	case SExpr:
		out += fmt.Sprintf("%s", l)
	case Atom:
		out += fmt.Sprintf("%s", l)
	}
	switch r := s.Right.(type) {
	case SExpr:
		out += fmt.Sprintf("%s", r.stringInner(false))
	case Atom:
		out += fmt.Sprintf(" . %s", r)
	}
	if top {
		out = out + ")"
	}
	return out
}

type Atom string

func (a Atom) isExpr() {}
func (a Atom) String() string {
	return string(a)
}

type Nil struct{}

func (n Nil) isExpr() {}
func (n Nil) String() string {
	return "NIL"
}

//tokens
type Token interface {
	isTok()
}

type LPAREN struct{}

func (l LPAREN) isTok() {}
func (l LPAREN) String() string {
	return "LPAREN"
}

type RPAREN struct{}

func (r RPAREN) isTok() {}
func (r RPAREN) String() string {
	return "RPAREN"
}

type DOT struct{}

func (d DOT) isTok() {}
func (d DOT) String() string {
	return "DOT"
}

type QUOTE struct{}

func (q QUOTE) isTok() {}
func (q QUOTE) String() string {
	return "QUOTE"
}

type NAME string

func (n NAME) isTok() {}

func Parse(tokens []Token) (Expr, int, error) {
	/*
		Rules are:
		NAME -> ATOM
		LPAREN RPAREN -> SExpr()
		LPAREN NAME RPAREN -> SExpr(ATOM)
		LPAREN NAME DOT NAME RPAREN -> SExpr(ATOM, ATOM)
		LPAREN NAME NAME RPAREN -> SExpr(ATOM, ATOM)
		LPAREN NAME NAME NAME RPAREN -> SExpr(ATOM, SExpr(ATOM, ATOM))

		LPAREN starts an SExpr
		RParen closes an SExpr
		NAME is an ATOM,
			if preceded by a LPAREN, is lval of SExpr
			if preceded by a NAME and in SExpr, then is

		SExpr with nothing it is valid:LPAREN RPAREN
		SExpr with only an Atom is valid: LPAREN NAME RPAREN
		SExpr with Atom DOT Atom is valid: LPAREN NAME DOT NAME RPAREN
		DOT must be followed by an LPAREN or NAME
	*/
	var out Expr
	rightSide := make([]Expr, 0)
	for k := 0; k < len(tokens); {
		v := tokens[k]
		switch t := v.(type) {
		case NAME:
			switch e := out.(type) {
			case nil:
				out = Atom(t)
			case SExpr:
				if e.Left == nil {
					e.Left = Atom(t)
				} else {
					rightSide = append(rightSide, Atom(t))
				}
				out = e
			case Atom:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))

			}
			k++
		case LPAREN:
			switch e := out.(type) {
			case nil:
				out = SExpr{}
				k++
			case SExpr:
				inner, pos, err := Parse(tokens[k:])
				if err != nil {
					return nil, pos, err
				}
				k = k + pos
				if e.Left == nil {
					e.Left = inner
				} else {
					rightSide = append(rightSide, inner)
				}
				out = e
			case Atom:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
			}
		case RPAREN:
			switch e := out.(type) {
			case nil:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
			case SExpr:
				//turn any rightside into a chain of s expressions
				l := len(rightSide)
				if l == 1 {
					e.Right = rightSide[0]
				} else if l > 1 {
					rightExpr := SExpr{rightSide[l-2], rightSide[l-1]}
					l = l - 3
					for ; l >= 0; l-- {
						rightExpr = SExpr{rightSide[l], rightExpr}
					}
					e.Right = rightExpr
				}
				return e, k + 1, nil
			case Atom:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
			}
		case DOT:
			switch e := out.(type) {
			case nil:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
			case SExpr:
				if e.Left == nil {
					return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
				}
				if len(rightSide) > 0 {
					return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
				}
			case Atom:
				return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ", t))
			}
			k++
		}
	}
	return out, len(tokens), nil
}

func Scan(s string) []Token {
	out := make([]Token, 0)
	curTokenTxt := make([]rune, 0)
	buildCurToken := func() {
		if len(curTokenTxt) > 0 {
			out = append(out, NAME(string(curTokenTxt)))
			curTokenTxt = make([]rune, 0)
		}
	}
	update := func(t Token) {
		buildCurToken()
		out = append(out, t)
	}
	for _, c := range s {
		switch {
		case c == '(':
			update(LPAREN{})
		case c == ')':
			update(RPAREN{})
		case c == '.':
			update(DOT{})
		case unicode.IsSpace(c):
			buildCurToken()
		case c == '\'':
			update(QUOTE{})
		default:
			curTokenTxt = append(curTokenTxt, c)
		}
	}
	buildCurToken()
	return out
}

func main() {
	fmt.Println(Scan("a"))
	fmt.Println(Scan("(a . b)"))
	fmt.Println(Scan("(a (b c))"))
	fmt.Println(Scan("(a b c)"))
	expr, pos, err := Parse(Scan("(a b c d)"))
	fmt.Println(expr)
	fmt.Println(pos)
	fmt.Println(err)
	bio := bufio.NewReader(os.Stdin)
	done := false
	depth := 0
	tokens := make([]Token, 0)
	for !done {
		line, _, err := bio.ReadLine()
		if err != nil {
			fmt.Errorf("Error: %v", err)
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
		line, hasMoreInLine, err := bio.ReadLine()
	}
}

package main

import (
	"fmt"
	"errors"
)

//expressions
type SExpr struct {
	Left Expr
	Right Expr
}

type Expr interface {
	isExpr()
}

func(s SExpr) isExpr() {}
func(s SExpr) String() string {
	out := "("
	hasLeft := false
	if s.Left != nil {
		out = out + fmt.Sprintf("%s",s.Left)
		hasLeft = true
	}
	if s.Right != nil {
		if hasLeft {
			out = out + " . "
		}
		out = out + fmt.Sprintf("%s", s.Right)
	}
	out = out + ")"
	return out
}

type Atom string

func(a Atom) isExpr() {}
func(a Atom) String() string {
	return string(a)
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
	rightSide := make([]Expr,0)
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
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
						
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
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
				}
			case RPAREN:
				switch e := out.(type) {
					case nil:
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
					case SExpr:
						//turn any rightside into a chain of s expressions
						l := len(rightSide)
						if l == 1 {
							e.Right = rightSide[0]
						} else if l > 1 {
							rightExpr := SExpr{rightSide[l-2], rightSide[l-1]}
							l = l -3
							for ; l >=0;l-- {
								rightExpr = SExpr{rightSide[l],rightExpr}
							}
							e.Right = rightExpr
						}
						return e, k+1, nil
					case Atom:
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
				}
			case DOT:
				switch e:= out.(type) {
					case nil:
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
					case SExpr:
						if e.Left == nil {
							return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
						}
						if len(rightSide) > 0 {
							return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
						}
					case Atom:
						return nil, k, errors.New(fmt.Sprintf("Illegal token: %s ",t))
				}
				k++
		}
	}
	return out, len(tokens), nil
}


func Scan(s string) []Token {
	out := make([]Token,0)
	curTokenTxt := make([]rune,0)
	for _,c := range s {
		switch c {
			case '(': {
			}
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,LPAREN{})
			case ')': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,RPAREN{})
			}
			case '.': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,DOT{})
			}
			case ' ': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
			}
			default: {
				curTokenTxt = append(curTokenTxt,c)
			}
		}
	}
	if len(curTokenTxt) > 0 {
		out = append(out, NAME(string(curTokenTxt)))
	}
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
}


package main

import (
	"fmt"
)

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
				return nil, k, fmt.Errorf("illegal token: %s ", t)

			}
			k++
		case LParen:
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
				return nil, k, fmt.Errorf("illegal token: %s ", t)
			}
		case RParen:
			switch e := out.(type) {
			case nil:
				return nil, k, fmt.Errorf("illegal token: %s ", t)
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
				return nil, k, fmt.Errorf("illegal token: %s ", t)
			}
		case Dot:
			switch e := out.(type) {
			case nil:
				return nil, k, fmt.Errorf("illegal token: %s ", t)
			case SExpr:
				if e.Left == nil {
					return nil, k, fmt.Errorf("illegal token: %s ", t)
				}
				if len(rightSide) > 0 {
					return nil, k, fmt.Errorf("illegal token: %s ", t)
				}
			case Atom:
				return nil, k, fmt.Errorf("illegal token: %s ", t)
			}
			k++
		}
	}
	return out, len(tokens), nil
}

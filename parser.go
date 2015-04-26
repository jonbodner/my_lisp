package main

import (
	"fmt"
)

func Parse(tokens []Token) (Expr, int, error) {
	return parseInner(tokens)
}

func parseInner(tokens []Token) (Expr, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("No tokens supplied.")
	}
	token := tokens[0]
	switch t := token.(type) {
	case NAME:
		//name by itself is a complete expression, so return
		out := Atom(t)
		return out, 1, nil
	case RParen:
		//this is an error
		return nil, 0, fmt.Errorf("Right Paren without matching Left Paren")
	case Dot:
		//this is an error
		return nil, 0, fmt.Errorf("Dot found not inside of List")
	case Quote:
		//"reader macro" -- turns 'EXPR into (QUOTE EXPR)
		quoted := SExpr{NIL, NIL}
		out := SExpr{Atom("QUOTE"), quoted}
		nested, remaining, error := parseInner(tokens[1:])
		if error != nil {
			return nil, remaining + 1, error
		}
		quoted.Left = nested
		return out, remaining + 1, nil
	case LParen:
		out := SExpr{NIL, NIL}
		cur := out
		dotted := false
		for k := 1; k < len(tokens); k++ {
			switch t2 := tokens[k].(type) {
			case NAME:
				newExpr := Atom(t2)
				error := buildSExpr(&cur, dotted, newExpr)
				if error != nil {
					return nil, k, error
				}
			case RParen:
				//done!
				return out, k + 1, nil
			case Dot:
				if !dotted {
					dotted = true
				} else {
					return nil, k, fmt.Errorf("Multiple Dots within a List")
				}
			case Quote, LParen:
				//recurse
				newExpr, remaining, error := parseInner(tokens[k:])
				k = k + remaining
				if error != nil {
					return nil, k, error
				}
				error = buildSExpr(&cur, dotted, newExpr)
				if error != nil {
					return nil, k, error
				}
			}
		}
		//fell off the end without finding RParen -- error!
		return nil, len(tokens), fmt.Errorf("Left Paren without matching Right Paren")
	}
	return nil, 0, fmt.Errorf("Unexpected Token found -- not processed!")
}

func buildSExpr(cur *SExpr, dotted bool, newExpr Expr) error {
	if cur.Left == NIL {
		cur.Left = newExpr
	}
	if cur.Right == NIL {
		if dotted {
			cur.Right = newExpr
		} else {
			cur.Right = SExpr{newExpr, NIL}
		}
	} else {
		if dotted {
			//error -- multiple items after the dot
			return fmt.Errorf("Multiple Expressions after a Dot in a List")
		} else {
			// (A (B NIL)) -> (A (B (C NIL)))
			newSExpr := SExpr{newExpr, NIL}
			cr := cur.Right.(SExpr)
			cr.Right = newSExpr
			fmt.Println("Added new link to chain", cur)
			*cur = cr
		}
	}
	fmt.Println("Should be pointing at last link in chain", cur)
	return nil
}

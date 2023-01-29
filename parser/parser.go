package parser

import "github.com/jonbodner/my_lisp/types"

func Parse(tokens []types.Token) (types.Expr, int, error) {
	return parseInner(tokens)
}

func parseInner(tokens []types.Token) (types.Expr, int, error) {
	//fmt.Println("incoming tokens:",tokens)
	if len(tokens) == 0 {
		return nil, 0, ParseError{"No tokens supplied", tokens, 0}
	}
	token := tokens[0]
	switch t := token.(type) {
	case types.NAME:
		//name by itself is a complete expression, so return
		out := types.Atom(t)
		return out, 1, nil
	case types.RParen:
		//this is an error
		return nil, 0, ParseError{"Right paren in unexpected location", tokens, 0}
	case types.Dot:
		//this is an error
		return nil, 0, ParseError{"Dot in unexpected location", tokens, 0}
	case types.Quote:
		//"reader macro" -- turns 'EXPR into (QUOTE EXPR)
		quoted := &types.SExpr{Left: types.NIL, Right: types.NIL}
		out := &types.SExpr{Left: types.Atom("QUOTE"), Right: quoted}
		nested, remaining, err := parseInner(tokens[1:])
		if err != nil {
			if pe, ok := err.(ParseError); ok {
				pe.pos += 1
				pe.tokens = tokens
				err = pe
			}
			return nil, remaining + 1, err
		}
		quoted.Left = nested
		return out, remaining + 1, nil
	case types.LParen:
		out := &types.SExpr{Left: types.NIL, Right: types.NIL}
		cur := out
		pos := 1
		dotted := false
		for {
			//fmt.Println("pos == ",pos)
			// if no more tokens, error
			if len(tokens) == pos {
				return nil, len(tokens), ParseError{"Left paren without matching right paren", tokens, 0}
			}
			// if the next token is RPAREN, we're done
			if tokens[pos] == types.RPAREN {
				return out, pos + 1, nil
			}
			//otherwise, recurse for the left value of the SExpr
			left, nextToken, err := parseInner(tokens[pos:])
			pos += nextToken
			if err != nil {
				if pe, ok := err.(ParseError); ok {
					pe.pos = pos
					pe.tokens = tokens
					err = pe
				}
				return nil, pos, err
			}
			//fmt.Println("got left value ",left, "to add to ", cur)
			cur.Left = left
			//if the next token is RPAREN, we're done
			if tokens[pos] == types.RPAREN {
				//fmt.Println("No right value -- done", out)
				return out, pos + 1, nil
			}
			//if the next token is a dot
			if tokens[pos] == types.DOT {
				if dotted {
					return nil, pos, ParseError{"More than one dot in a dotted pair", tokens, pos}
				}
				dotted = true
				pos++
				right, nextToken, err := parseInner(tokens[pos:])
				pos += nextToken
				if err != nil {
					if pe, ok := err.(ParseError); ok {
						pe.pos = pos
						pe.tokens = tokens
						err = pe
					}
					return nil, pos, err
				}
				cur.Right = right
			} else {
				if dotted {
					return nil, pos, ParseError{"More than one value to the right of the dot in a dotted pair", tokens, pos}
				}
				//otherwise, keep going
				right := &types.SExpr{Left: types.NIL, Right: types.NIL}
				cur.Right = right
				cur = right
			}
		}

	}
	return nil, 0, ParseError{"Unexpected Token found -- not processed!", tokens, 0}
}

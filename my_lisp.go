package main

import (
	"fmt"
	"errors"
	"bufio"
	"os"
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

/*
1. value (QUOTE e) = e. Thus the value of (QUOTE A) is A.
2. value (CAR e), where value e is a non-empty list, is the first
element of value e. Thus value (CAR (QUOTE (A B C))) - A.
3. value (CDR e), where value e is a non-empty list, is the the list
that remains when the first element of value e is deleted. Thus
value (CDR (QUOTE (A B C))) = (B C).
4. value (CONS el e2), is the list that results from prefixing
value el onto the list value e2. Thus
value (CONS (QUOTE A) (QUOTE (B C))) = (A B C).
5. value (EQUAL el e2) is T if value el = value e2. Otherwise, its
value is NIL. Thus
value (EQUAL (CAR (QUOTE (A B))) (QUOTE A)) = T,
6. value (ATOM e) - T if value e is an atom; otherwise its value Is
NIL.
7. value (COND(pt e I) ... (PB en)) = value e i, where Pi is the the
first of the p's whose value is not NIL. Thus
value (COND ((ATOM (QUOTE A)) (QUOTE B)) ((QUOTE T)
(QUOTE C))) = B.
8. An atom v, regarded as a variable, may have a value.
9. value ((LAMBDA (v I ... v,) e) e I ... e a) is the same as value e
but in an environment in which the variables v I ... v n take the
values of the expressions e I ... e n in the original environment.
Thus
value ((LAMBDA (X Y) (CONS (CAR X) Y)) (QUOTE (A B))
(CDR (QUOTE (C D)))) = (A D).
10. Here's the hard one. value ((LABEL f (LAMBDA (o I ... v,)
e)) e I ... en) is the same as value ((LAMBDA (v I ... vn) e) e I ...
e n) with the additional rule that whenever ~¢al an) must be
evaluated, f is replaced by (LABEL/" (LAMBDA (v I ... vn) e)).
Lists beginning with LABEL define functions recursively.
 */
func Eval(e Expr) (Expr, error) {
	switch t := e.(type) {
	case Atom:
		//look up variable value in context and return that
		return nil, errors.New("Not implemented yet")
	case SExpr:
		switch a := t.Left.(type) {
		case Atom:
			switch a {
			case "quote":
				if t.Right == nil {
					return nil, errors.New("Missing parameter for quote")
				}
				return t.Right, nil
			case "car":
				e2, error := Eval(t.Right)
				if error != nil {
					return nil, error
				}
				switch a2 := e2.(type) {
				case Atom:
					return nil, errors.New("car parameter must be a list")
				case SExpr:
					return a2.Left, nil
				}
			case "cdr":
				e2, error := Eval(t.Right)
				if error != nil {
					return nil, error
				}
				switch a2 := e2.(type) {
				case Atom:
					return nil, errors.New("cdr parameter must be a list")
				case SExpr:
					return a2.Right, nil
				}
			case "cons":
				switch a2 := t.Right.(type) {
				case Atom:
					return nil, errors.New("cons needs two parameters")
				case SExpr:
					e1, error := Eval(a2.Left)
					if error != nil {
						return nil, error
					}
					e2, error := Eval(a2.Right)
					if error != nil {
						return nil, error
					}
					return SExpr{e1, e2}, nil

				}
			case "equal":
				return nil, errors.New("Not implemented yet")
			case "cond":
				return nil, errors.New("Not implemented yet")
			case "lambda":
				return nil, errors.New("Not implemented yet")
			case "label":
			default:
				//look up variable value in context and return that
				return nil, errors.New("Not implemented yet")
			}
		case SExpr:
			return nil, errors.New("No function specified")
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("Don't know how I got here")
}

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


func Scan(s string) ([]Token, int) {
	out := make([]Token,0)
	curTokenTxt := make([]rune,0)
	depth := 0
	for _,c := range s {
		switch c {
			case '(': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,LPAREN{})
				depth++
			}
			case ')': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,RPAREN{})
				depth--
			}
			case '.': {
				if len(curTokenTxt) > 0 {
					out = append(out, NAME(string(curTokenTxt)))
					curTokenTxt = make([]rune,0)
				}
				out = append(out,DOT{})
			}
			case '\n', '\r', '\t', ' ': {
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
	return out, depth
}

func main() {
	bio := bufio.NewReader(os.Stdin)
	done := false
	depth := 0
	tokens := make([]Token,0)
	for !done {
		line, _ , err := bio.ReadLine()
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
			tokens = make([]Token,0)
		}
	}
}


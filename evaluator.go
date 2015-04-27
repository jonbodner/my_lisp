package main

import (
	"errors"
	"fmt"
)

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
	fmt.Println("Evaluating ", e)
	switch t := e.(type) {
	case Atom:
		//look up variable value in context and return that
		return nil, errors.New("not implemented yet")
	case *SExpr:
		switch a := t.Left.(type) {
		case Atom:
			switch a {
			case "QUOTE":
				if t.Right == NIL {
					return nil, errors.New("missing parameter for QUOTE")
				}
				switch a2 := t.Right.(type) {
					case Atom:
						return nil, errors.New("shouldn't have an Atom after a QUOTE")
					case *SExpr:
						//should only have a single parameter for QUOTE
						if a2.Right != NIL {
							return nil, errors.New("shouldn't have more than one parameter for QUOTE")
						}
						return a2.Left, nil
				}
			case "CAR":
				if t.Right == NIL {
					return nil, errors.New("missing parameter for CAR")
				}
				switch a2 := t.Right.(type) {
					case Atom:
						return nil, errors.New("CAR parameter must be a list")
					case *SExpr:
						//should only have a single parameter for CAR
						if a2.Right != NIL {
							return nil, errors.New("shouldn't have more than one parameter for CAR")
						}
						e2, error := Eval(a2.Left)
						if error != nil {
							return nil, error
						}
						switch a3 := e2.(type) {
							case Atom:
								return nil, errors.New("CAR parameter must be a list")
							case *SExpr:
								return a3.Left, nil
						}
				}
			case "CDR":
				if t.Right == NIL {
					return nil, errors.New("missing parameter for CDR")
				}
				switch a2 := t.Right.(type) {
					case Atom:
					return nil, errors.New("CDR parameter must be a list")
					case *SExpr:
					//should only have a single parameter for CDR
					if a2.Right != NIL {
						return nil, errors.New("shouldn't have more than one parameter for CDR")
					}
					e2, error := Eval(a2.Left)
					if error != nil {
						return nil, error
					}
					switch a3 := e2.(type) {
						case Atom:
						return nil, errors.New("CDR parameter must be a list")
						case *SExpr:
						return a3.Right, nil
					}
				}
			case "CONS":
				//must have two params
				//going to construct an SExpr out of them
				//first is going to be the left, second is going to be the right
				if t.Right == NIL {
					return nil, errors.New("missing parameters for CONS")
				}
				switch a2 := t.Right.(type) {
					case Atom:
					return nil, errors.New("CONS parameter must be a list")
					case *SExpr:
					e2, error := Eval(a2.Left)
					if error != nil {
						return nil, error
					}
					//should have two parameters for CDR
					if a2.Right == NIL {
						return nil, errors.New("must have two parameters for CONS")
					}
					switch a3 := a2.Right.(type) {
						case Atom:
						return nil, errors.New("CONS parameter must be a list")
						case *SExpr:
						if a3.Right != NIL {
							return nil, errors.New("must have two parameters for CONS")
						}
						e3, error := Eval(a3.Left)
						if error != nil {
							return nil, error
						}
						return &SExpr{e2, e3}, nil
					}
				}
			case "ATOM":
				if t.Right == NIL {
					return nil, errors.New("missing parameter for ATOM")
				}
				switch a2 := t.Right.(type) {
					case Atom:
					return Atom("T"), nil
					case *SExpr:
					//should only have a single parameter for ATOM
					if a2.Right != NIL {
						return nil, errors.New("shouldn't have more than one parameter for ATOM")
					}
					e2, error := Eval(a2.Left)
					if error != nil {
						return nil, error
					}
					switch a3 := e2.(type) {
						case Atom:
						return Atom("T"), nil
						case *SExpr:
						if a3.Left == NIL && a3.Right == NIL {
							return Atom("T"), nil
						}
						return &SExpr{NIL, NIL}, nil
					}
				}
			case "EQUAL":
			//must have two params
			if t.Right == NIL {
				return nil, errors.New("missing parameters for EQUAL")
			}
			switch a2 := t.Right.(type) {
				case Atom:
				return nil, errors.New("EQUAL parameter must be a list")
				case *SExpr:
				e2, error := Eval(a2.Left)
				if error != nil {
					return nil, error
				}
				//should have two parameters for EQUAL
				if a2.Right == NIL {
					return nil, errors.New("must have two parameters for EQUAL")
				}
				switch a3 := a2.Right.(type) {
					case Atom:
					return nil, errors.New("EQUAL parameter must be a list")
					case *SExpr:
					if a3.Right != NIL {
						return nil, errors.New("must have two parameters for EQUAL")
					}
					e3, error := Eval(a3.Left)
					if error != nil {
						return nil, error
					}
					if isEqual(e2, e3) {
						return Atom("T"), nil
					}
					return &SExpr{NIL, NIL}, nil
				}
			}
			case "cond":
				return nil, errors.New("not implemented yet")
			case "lambda":
				return nil, errors.New("not implemented yet")
			case "label":
			default:
				//look up variable value in context and return that
				return nil, errors.New("not implemented yet")
			}
		case *SExpr:
			return nil, errors.New("no function specified")
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("don't know how I got here")
}

func isEqual(e, e2 Expr) bool {
	switch e := e.(type) {
		case Atom:
		if e2, ok := e2.(Atom); ok {
			return e == e2
		}
		return false
		case Nil:
		_, ok := e2.(Nil)
		return ok
		case *SExpr:
		if e2, ok := e2.(*SExpr); ok {
			return isEqual(e.Left, e2.Left) && isEqual(e.Right, e2.Right)
		}
		return false
	}
	return false
}
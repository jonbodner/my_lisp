package main

import (
	"errors"
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
	switch t := e.(type) {
	case Atom:
		//look up variable value in context and return that
		return nil, errors.New("not implemented yet")
	case SExpr:
		switch a := t.Left.(type) {
		case Atom:
			switch a {
			case "quote":
				if t.Right == nil {
					return nil, errors.New("missing parameter for quote")
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
				return nil, errors.New("not implemented yet")
			case "cond":
				return nil, errors.New("not implemented yet")
			case "lambda":
				return nil, errors.New("not implemented yet")
			case "label":
			default:
				//look up variable value in context and return that
				return nil, errors.New("not implemented yet")
			}
		case SExpr:
			return nil, errors.New("no function specified")
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("don't know how I got here")
}

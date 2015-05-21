package main

import (
	"errors"
	"fmt"
	"math/big"
)

/*
DONE 1. value (QUOTE e) = e. Thus the value of (QUOTE A) is A.
DONE 2. value (CAR e), where value e is a non-empty list, is the first
element of value e. Thus value (CAR (QUOTE (A B C))) - A.
DONE 3. value (CDR e), where value e is a non-empty list, is the the list
that remains when the first element of value e is deleted. Thus
value (CDR (QUOTE (A B C))) = (B C).
DONE 4. value (CONS el e2), is the list that results from prefixing
value el onto the list value e2. Thus
value (CONS (QUOTE A) (QUOTE (B C))) = (A B C).
DONE 5. value (EQUAL el e2) is T if value el = value e2. Otherwise, its
value is NIL. Thus
value (EQUAL (CAR (QUOTE (A B))) (QUOTE A)) = T,
DONE 6. value (ATOM e) - T if value e is an atom; otherwise its value Is
NIL.
DONE 7. value (COND(pt e I) ... (PB en)) = value e i, where Pi is the the
first of the p's whose value is not NIL. Thus
value (COND ((ATOM (QUOTE A)) (QUOTE B)) ((QUOTE T)
(QUOTE C))) = B.
DONE 8. An atom v, regarded as a variable, may have a value.
9. value ((LAMBDA (v I ... v,) e) e I ... e a) is the same as value e
but in an environment in which the variables v I ... v n take the
values of the expressions e I ... e n in the original environment.
Thus
value ((LAMBDA (X Y) (CONS (CAR X) Y)) (QUOTE (A B))
(CDR (QUOTE (C D)))) = (A D).
DONE 10. Here's the hard one. value ((LABEL f (LAMBDA (o I ... v,)
e)) e I ... en) is the same as value ((LAMBDA (v I ... vn) e) e I ...
e n) with the additional rule that whenever ~¢al an) must be
evaluated, f is replaced by (LABEL/" (LAMBDA (v I ... vn) e)).
Lists beginning with LABEL define functions recursively.
*/

type Env map[Atom]Expr

type Evaluator func(*SExpr, Env) (Expr, error)

var GlobalEnv = make(Env)

var BuiltIn map[Atom]Evaluator

func init() {
	GlobalEnv[T] = T
	GlobalEnv[Atom("NIL")] = EMPTY

	BuiltIn = map[Atom]Evaluator{
		Atom("QUOTE"): quote,
		Atom("CAR"):   car,
		Atom("CDR"):   cdr,
		Atom("CONS"):  cons,
		Atom("ATOM"):  atom,
		Atom("EQ"):    equal,
		Atom("COND"):  cond,
		Atom("LABEL"): label,
	}
}

func Eval(e Expr) (Expr, error) {
	return evalInner(e, GlobalEnv)
}

func evalInner(e Expr, env Env) (Expr, error) {
	fmt.Println("Evaluating ", e)
	switch t := e.(type) {
	case Atom:
		//look up variable value in context and return that
		expr, ok := env[t]
		if ok {
			if _, ok = expr.(Atom); ok {
				return expr, nil
			}
			return evalInner(expr, env)
		}
		//check if number, and if so return self
		r := &big.Rat{}
		_, ok = r.SetString(string(t))
		if ok {
			return t, nil
		}
		return nil, fmt.Errorf("Unknown symbol %s ", t)
	case *SExpr:
		switch a := t.Left.(type) {
		case Atom:
			evaluator, ok := BuiltIn[a]
			if ok {
				return evaluator(t, env)
			}
			//look up variable value in context and process that
			fmt.Println("looking up ", a)
			expr, ok := env[a]
			if !ok {
				return nil, fmt.Errorf("Unknown symbol %s ", a)
			}
			//replace the atom with the value of the expression
			result, err := evalInner(expr, env)
			if err != nil {
				return nil, err
			}
			t.Left = result
			return evalInner(t, env)
		case *SExpr:
			return nil, errors.New("no function specified")
		case Nil:
			return t, nil
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("don't know how I got here")
}

func quote(t *SExpr, env Env) (Expr, error) {
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
	default:
		return nil, errors.New("shouldn't get here")
	}
}

func car(t *SExpr, env Env) (Expr, error) {
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
		e2, error := evalInner(a2.Left, env)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			return nil, errors.New("CAR parameter must be a list")
		case *SExpr:
			return a3.Left, nil
		}
	default:
		return nil, fmt.Errorf("Unknown Expr type found: %T", a2)
	}
	return nil, errors.New("Should never get here")
}

func cdr(t *SExpr, env Env) (Expr, error) {
	fmt.Println("passed in: ", t)
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
		e2, error := evalInner(a2.Left, env)
		fmt.Println("after eval, e2 == ", e2)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			return nil, errors.New("CDR parameter must be a list")
		case *SExpr:
			fmt.Println("returning ", a3.Right)
			return a3.Right, nil
		}
	default:
		return nil, fmt.Errorf("Unknown Expr type found: %T", a2)
	}
	return nil, errors.New("Should never get here")
}

func cons(t *SExpr, env Env) (Expr, error) {
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
		e2, error := evalInner(a2.Left, env)
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
			e3, error := evalInner(a3.Left, env)
			if error != nil {
				return nil, error
			}
			return &SExpr{e2, e3}, nil
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("Shouldn't get here")
}

func atom(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameter for ATOM")
	}
	switch a2 := t.Right.(type) {
	case Atom:
		return T, nil
	case *SExpr:
		//should only have a single parameter for ATOM
		if a2.Right != NIL {
			return nil, errors.New("shouldn't have more than one parameter for ATOM")
		}
		e2, error := evalInner(a2.Left, env)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			return T, nil
		case *SExpr:
			if a3.Left == NIL && a3.Right == NIL {
				return T, nil
			}
			return EMPTY, nil
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("Shouldn't get here")
}

func equal(t *SExpr, env Env) (Expr, error) {
	//must have two params
	if t.Right == NIL {
		return nil, errors.New("missing parameters for EQUAL")
	}
	switch a2 := t.Right.(type) {
	case Atom:
		return nil, errors.New("EQUAL parameter must be a list")
	case *SExpr:
		e2, error := evalInner(a2.Left, env)
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
			e3, error := evalInner(a3.Left, env)
			if error != nil {
				return nil, error
			}
			if isEqual(e2, e3) {
				return T, nil
			}
			return EMPTY, nil
		}
	default:
		return nil, errors.New("shouldn't get here")
	}
	return nil, errors.New("Shouldn't get here")
}

func cond(t *SExpr, env Env) (Expr, error) {
	//find the first non-NIL result, and return it
	pos := 1
	for {
		cur, error := nth(pos, t)
		fmt.Println(pos, " item is ", cur)
		if error != nil {
			return nil, error
		}
		switch cur := cur.(type) {
		case Atom:
			return nil, errors.New("Cannot have an atom as a COND parameter")
		case Nil:
			return EMPTY, nil
		case *SExpr:
			fmt.Println("Evaluating ", cur.Left)
			car, error := evalInner(cur.Left, env)
			fmt.Println("Response: ", car, error)
			if error != nil {
				return nil, error
			}
			if !isEqual(car, EMPTY) {
				switch result := cur.Right.(type) {
				case Atom:
					return nil, errors.New("Cannot have a dotted pair here")
				case Nil:
					return EMPTY, nil
				case *SExpr:
					return evalInner(result.Left, env)
				}
			}
		}
		pos++
	}
}

func label(t *SExpr, env Env) (Expr, error) {
	//must have two params
	//first must be an atom
	//second can be any expression
	//going to assign expression to atom
	if t.Right == NIL {
		return nil, errors.New("missing parameters for LABEL")
	}
	switch a2 := t.Right.(type) {
	case Atom:
		return nil, errors.New("LABEL parameter must be a list")
	case *SExpr:
		//a2.Left must be an atom
		l, ok := a2.Left.(Atom)
		if !ok {
			return nil, errors.New("LABEL can only be assigned to an Atom")
		}
		//a2.Right must be an *SExpr
		a3, ok := a2.Right.(*SExpr)
		if !ok {
			return nil, errors.New("LABEL parameter must be a list")
		}
		//a2.Right.Right must be NIL
		if a3.Right != NIL {
			return nil, errors.New("must have two parameters for LABEL")
		}
		//a2.Right.Left can be anything
		lval := a3.Left
		env[l] = lval
		return l, nil
	}
	return nil, errors.New("Shouldn't get here")
}

//get the nth parameter of the SExpr.
//The function/macro/special form name is the CAR of the SExpr passed in
//pos == 1 for the first parameter. This is the CAR of the CDR of the SExpr passed in
//If there are not enough parameters, NIL is returned
func nth(pos int, e *SExpr) (Expr, error) {
	for i := 0; i < pos; i++ {
		fmt.Println(i, "e is ", e)
		next := e.Right
		switch next := next.(type) {
		case Atom:
			return nil, errors.New("Can't have a dotted pair here")
		case *SExpr:
			e = next
		case Nil:
			return NIL, nil
		}
	}
	return e.Left, nil
}

func isEqual(e, e2 Expr) bool {
	switch e := e.(type) {
	case Atom:
		if e2, ok := e2.(Atom); ok {
			if e == e2 {
				return true
			} else {
				//might be numbers, have to do numeric comparison
				r := &big.Rat{}
				r, ok := r.SetString(string(e))
				r2 := &big.Rat{}
				r2, ok2 := r2.SetString(string(e2))
				if ok && ok2 {
					return r.Cmp(r2) == 0
				}
			}
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

/*
todo

macros
slices
maps
channels
go routines
select
numeric operations
load/save environment
*/

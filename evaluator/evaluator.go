package evaluator

import (
	"errors"
	"fmt"
	"math/big"
	. "host.bodnerfamily.com/my_lisp/types"
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
DONE 9. value ((LAMBDA (v I ... v,) e) e I ... e a) is the same as value e
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

type Evaluator func(*SExpr, Env) (Expr, error)

var TopLevel = make(GlobalEnv)

var BuiltIn map[Atom]Evaluator

func init() {
	TopLevel[T] = T
	TopLevel[Atom("NIL")] = EMPTY

	BuiltIn = map[Atom]Evaluator{
		Atom("QUOTE"): quote,
		Atom("CAR"):   car,
		Atom("CDR"):   cdr,
		Atom("CONS"):  cons,
		Atom("ATOM"):  atom,
		Atom("EQ"):    equal,
		Atom("COND"):  cond,
		Atom("LABEL"): label,
		Atom("SETQ"):  setq,
		Atom("LAMBDA"):  lambda,
		Atom("+"): plus,
		Atom("-"): minus,
		Atom("*"): times,
		Atom("/"): div,
	}
}

func Eval(e Expr) (Expr, error) {
	return evalInner(e, TopLevel)
}

var depth = 0

func evalInner(e Expr, env Env) (Expr, error) {
	depth++
	defer func () {
		fmt.Println("done depth ", depth)
		depth--
	}()
	fmt.Println("at depth ", depth)
	fmt.Println("Evaluating ", e)
	switch t := e.(type) {
	case Atom:
		fmt.Println("\tGot an Atom")
		//check if number, and if so return self
		r := &big.Rat{}
		_, ok := r.SetString(string(t))
		if ok {
			fmt.Println("\t\t\tand it's a number so we're done")
			return t, nil
		}
		//look up variable value in context and return that
		expr, ok := env.Get(t)
		if ok {
			fmt.Println("\t\tvalue found for atom")
			return expr, nil
		}
		return nil, fmt.Errorf("Unknown symbol %s ", t)
	case *SExpr:
		fmt.Println("\tGot an SExpr")
		switch a := t.Left.(type) {
		case Atom:
			fmt.Println("\t\tLeft is an Atom")
			evaluator, ok := BuiltIn[a]
			if ok {
				fmt.Println("running ", a, evaluator)
				return evaluator(t, env)
			}
			fmt.Println("not a builtin")
			//look up variable value in context and process that
			fmt.Println("looking up ", a)
			expr, ok := env.Get(a)
			if !ok {
				return nil, fmt.Errorf("Unknown symbol %s ", a)
			}
			//replace the atom with the value of the expression
			result, err := evalInner(expr, env)
			fmt.Println("done evaluating")
			if err != nil {
				return nil, err
			}
			newT := &SExpr{result, t.Right}
			return evalInner(newT, env)
		case *SExpr:
			fmt.Println("\t\tLeft is an SExpr")
			//evaluate the left, then replace left with the evaluated value, and recurse
			lResult, err := evalInner(t.Left, env)
			if err != nil {
				return nil, err
			}
			t.Left = lResult
			return evalInner(t, env)
		case Nil:
			fmt.Println("Got a nil left")
			return t, nil
		case Lambda:
			fmt.Println("\t\tLeft is a Lambda")
			return processLambda(a, t, env)
		default:
			return nil, errors.New("shouldn't get here")
		}
	case Lambda:
		fmt.Println("\tGot a lambda")
		return t, nil
	default:
		fmt.Println("\t what happened?")
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
			fmt.Println("second cons param:", e3)
			if e3 == EMPTY {
				e3 = NIL
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
		env.Put(l, lval)
		return l, nil
	}
	return nil, errors.New("Shouldn't get here")
}

func setq(t *SExpr, env Env) (Expr, error) {
	//must have two params
	//first must be an atom
	//second can be any expression
	//going to evaluate expression and assign result to atom
	if t.Right == NIL {
		return nil, errors.New("missing parameters for SETQ")
	}
	switch a2 := t.Right.(type) {
	case Atom:
		return nil, errors.New("SETQ parameter must be a list")
	case *SExpr:
		//a2.Left must be an atom
		l, ok := a2.Left.(Atom)
		if !ok {
			return nil, errors.New("SETQ can only be assigned to an Atom")
		}
		//a2.Right must be an *SExpr
		a3, ok := a2.Right.(*SExpr)
		if !ok {
			return nil, errors.New("SETQ parameter must be a list")
		}
		//a2.Right.Right must be NIL
		if a3.Right != NIL {
			return nil, errors.New("must have two parameters for SETQ")
		}
		//a2.Right.Left can be anything
		lval, err := evalInner(a3.Left, env)
		if err != nil {
			return nil, err
		}
		env.Put(l, lval)
		return l, nil
	}
	return nil, errors.New("Shouldn't get here")
}

func plus(t *SExpr, env Env) (Expr, error) {
	fmt.Println("in +")
	if t.Right == NIL {
		return nil, errors.New("missing parameters for +")
	}
	params, ok := t.Right.(*SExpr)
	if !ok {
		return nil, errors.New("+ parameters must be a list")
	}
	r := &big.Rat{}
	for {
		ev, err := evalInner(params.Left, env)
		if err != nil {
			return nil, err
		}
		fmt.Println("\tfinished eval -- checking if it's a number")
		r2 := &big.Rat{}
		_, ok = r2.SetString(ev.String())
		if !ok {
			return nil, fmt.Errorf("%s is not a valid number", ev)
		}
		r.Add(r, r2)
		next := params.Right
		if next == NIL {
			break
		}
		n, ok := next.(*SExpr)
		if !ok {
			return nil, errors.New("+ parameters must be a List")
		}
		params = n
	}
	return Atom(r.RatString()), nil
}

func minus(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameters for -")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == NIL {
		return nil, errors.New("- requires at least one parameter")
	}
	r := &big.Rat{}
	first := true
	for {
		ev, err := evalInner(v, env)
		if err != nil {
			return nil, err
		}
		r2 := &big.Rat{}
		_, ok := r2.SetString(ev.String())
		if !ok {
			return nil, fmt.Errorf("%s is not a valid number", ev)
		}
		if first {
			r = r2
		} else {
			r.Sub(r, r2)
		}
		pos++
		v, err = nth(pos, t)
		if err != nil {
			return nil, err
		}
		if v == NIL {
			//if there was only one value, just negate it
			if first {
				r.Neg(r)
			}
			break
		}
		first = false
	}
	return Atom(r.RatString()), nil
}

func times(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameters for *")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == NIL {
		return nil, errors.New("- requires at least one parameter")
	}
	r := &big.Rat{}
	r.SetString("1")
	for {
		ev, err := evalInner(v, env)
		if err != nil {
			return nil, err
		}
		r2 := &big.Rat{}
		_, ok := r2.SetString(ev.String())
		if !ok {
			return nil, fmt.Errorf("%s is not a valid number", ev)
		}
		r.Mul(r, r2)
		pos++
		v, err = nth(pos, t)
		if err != nil {
			return nil, err
		}
		if v == NIL {
			break
		}
	}
	return Atom(r.RatString()), nil
}

func div(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameters for /")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == NIL {
		return nil, errors.New("/ requires at least one parameter")
	}
	r := &big.Rat{}
	first := true
	for {
		ev, err := evalInner(v, env)
		if err != nil {
			return nil, err
		}
		r2 := &big.Rat{}
		_, ok := r2.SetString(ev.String())
		if !ok {
			return nil, fmt.Errorf("%s is not a valid number", ev)
		}
		if first {
			r = r2
		} else {
			r.Mul(r, r2.Inv(r2))
		}
		pos++
		v, err = nth(pos, t)
		if err != nil {
			return nil, err
		}
		if v == NIL {
			//if there was only one value, just negate it
			if first {
				r.Inv(r)
			}
			break
		}
		first = false
	}
	return Atom(r.RatString()), nil
}

func lambda(t *SExpr, env Env) (Expr, error) {
	//must have 2 params
	//param 1 is a list of parameters
	params, err := nth(1, t)
	if err != nil {
		return nil, err
	}
	if params == NIL {
		return nil, errors.New("missing parameters for LAMBDA")
	}
	l, ok := params.(*SExpr)
	if !ok {
		return nil, errors.New("LAMBDA parameter list must be a List")
	}

	//copy into slice of Atoms
	aList, err := listToSlice(l)
	if err != nil {
		return nil, err
	}

	//param 2 is an Expr
	body, err := nth(2, t)
	if err != nil {
		return nil, err
	}
	if body == NIL {
		return nil, errors.New("must have two parameters for LAMBDA")
	}
	//a2.Right.Left can be anything
	//returns a new Expr type, a Lambda, which has its own env
	lambda := Lambda{ParentEnv:env, Body: body, Params: aList}
	return lambda, nil
}

func listToSlice(l *SExpr) ([]Atom, error) {
	out := []Atom{}

	pos := 0
	for {
		cur, err := nth(pos, l)
		if err != nil {
			return nil, err
		}

		if cur == NIL {
			break
		}

		c, ok := cur.(Atom)
		if !ok {
			return nil, errors.New("Only Atoms can be parameter names")
		}

		out = append(out, c)
		pos++
	}
	return out, nil
}

func processLambda(l Lambda, t *SExpr, env Env) (Expr, error) {
	fmt.Println("in processLambda")
	le := LocalEnv{Vals:make(map[Atom]Expr), Parent: l.ParentEnv}
	//assign parameter values to parameter names
	count :=0
	switch paramVals := t.Right.(type) {
		case Atom:
			return nil, errors.New("Can't have a dotted pair here")
		case Nil:
			//do nothing
		case *SExpr:
			var param Expr = NIL
			for count < len(l.Params) {
				var err error
				param, err = nth(count, paramVals)
				if err != nil {
					return nil, err
				}
				if param == NIL {
					break
				}
				val, err := evalInner(param, env)
				if err != nil {
					return nil, err
				}
				le.Vals[l.Params[count]] = val
				count++
			}
			leftOver, err := nth(len(l.Params), paramVals)
			if err != nil {
				return nil, err
			}
			if leftOver != NIL {
				return nil, fmt.Errorf("Too many parameters for LAMBDA. Expected %d", len(l.Params))
			}
	}
	if count != len(l.Params) {
		return nil, fmt.Errorf("Too few parameters for LAMBDA. Expected %d, got %d", len(l.Params), count)
	}
	fmt.Println("\tevaling body with envronment", le)
	//call body with new environment
	return evalInner(l.Body, le)
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
let/scopes
progn
load/save environment
numeric operations
code cleanup
logging levels
macros

experiments:
slices
maps
channels
go routines
select
tail call optimization?
compilation?
*/

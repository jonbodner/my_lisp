package evaluator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/jonbodner/my_lisp/global"
	"github.com/jonbodner/my_lisp/parser"
	"github.com/jonbodner/my_lisp/scanner"
	. "github.com/jonbodner/my_lisp/types"
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
		Atom("QUOTE"):     quote,
		Atom("CAR"):       car,
		Atom("CDR"):       cdr,
		Atom("CONS"):      cons,
		Atom("ATOM"):      atom,
		Atom("EQ"):        equal,
		Atom("COND"):      cond,
		Atom("LABEL"):     label,
		Atom("SETQ"):      setq,
		Atom("LAMBDA"):    lambda,
		Atom("PROGN"):     progn,
		Atom("LET"):       let,
		Atom("**DEBUG**"): debug,
		Atom("LOAD"):      load,
		Atom("STORE"):     store,
		Atom("DELETE"):    delete,
	}
}

func Eval(e Expr) (Expr, error) {
	return evalInner(e, TopLevel)
}

var depth = 0

func evalInner(e Expr, env Env) (Expr, error) {
	depth++
	defer func() {
		global.Log("done depth ", depth)
		depth--
	}()
	global.Log("at depth ", depth)
	global.Log("Evaluating ", e)
	switch t := e.(type) {
	case Atom:
		global.Log("\tGot an Atom")
		//check if number, and if so return self
		r := &big.Rat{}
		_, ok := r.SetString(string(t))
		if ok {
			return t, nil
		}
		//look up variable value in context and return that
		expr, ok := env.Get(t)
		if ok {
			return expr, nil
		}
		return nil, fmt.Errorf("Unknown symbol %s ", t)
	case *SExpr:
		global.Log("\tGot an SExpr")
		switch a := t.Left.(type) {
		case Atom:
			global.Log("\t\tLeft is an Atom")
			evaluator, ok := BuiltIn[a]
			if ok {
				return evaluator(t, env)
			}
			global.Log("not a builtin")
			//look up variable value in context and process that
			global.Log("looking up ", a)
			expr, ok := env.Get(a)
			if !ok {
				return nil, fmt.Errorf("Unknown symbol %s ", a)
			}
			//replace the atom with the value of the expression
			result, err := evalInner(expr, env)
			global.Log("done evaluating")
			if err != nil {
				return nil, err
			}
			newT := &SExpr{result, t.Right}
			return evalInner(newT, env)
		case *SExpr:
			global.Log("\t\tLeft is an SExpr")
			//evaluate the left, then replace left with the evaluated value, and recurse
			lResult, err := evalInner(t.Left, env)
			if err != nil {
				return nil, err
			}
			t.Left = lResult
			return evalInner(t, env)
		case Nil:
			global.Log("Got a nil left")
			return t, nil
		case Lambda:
			global.Log("\t\tLeft is a Lambda")
			return processLambda(a, t, env)
		default:
			return nil, errors.New("shouldn't get here")
		}
	case Lambda:
		global.Log("\tGot a lambda")
		return t, nil
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
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			return nil, errors.New("CDR parameter must be a list")
		case *SExpr:
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
		if error != nil {
			return nil, error
		}
		switch cur := cur.(type) {
		case Atom:
			return nil, errors.New("Cannot have an atom as a COND parameter")
		case Nil:
			return EMPTY, nil
		case *SExpr:
			car, error := evalInner(cur.Left, env)
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
		return lval, nil
	}
	return nil, errors.New("Shouldn't get here")
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
	lambda := Lambda{ParentEnv: env, Body: body, Params: aList}
	return lambda, nil
}

//load the environment from the named file. Existing symbols will be overwritten.
func load(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameter for LOAD")
	}
	switch a2 := t.Right.(type) {
	case *SExpr:
		//should only have a single parameter for LOAD
		if a2.Right != NIL {
			return nil, errors.New("shouldn't have more than one parameter for LOAD")
		}
		e2, error := evalInner(a2.Left, env)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			f, err := os.Open(string(a3))
			if err != nil {
				return nil, err
			}
			defer f.Close()
			newEnv, err := internalRepl(f)
			if err != nil {
				return nil, err
			}
			for k, v := range newEnv {
				TopLevel[k] = v
			}
			return T, nil
		default:
			return nil, errors.New("LOAD parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("Shouldn't get here")
}

func internalRepl(r io.Reader) (GlobalEnv, error) {
	newEnv := GlobalEnv{}
	newEnv[T] = T
	newEnv[Atom("NIL")] = EMPTY

	bio := bufio.NewReader(r)
	done := false
	depth := 0
	tokens := []Token{}
	for !done {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			done = true
			continue
		}
		newTokens, newDepth := scanner.Scan(string(line))
		depth = depth + newDepth
		if depth < 0 {
			return nil, errors.New("Invalid -- Too many closing parens")
		}
		tokens = append(tokens, newTokens...)
		if depth == 0 {
			//global.Log(tokens)
			expr, _, err := parser.Parse(tokens)
			//global.Log(expr)
			//global.Log(pos)
			//global.Log(err)
			if err != nil {
				global.Log(err)
			} else {
				result, err := evalInner(expr, newEnv)
				if err != nil {
					return nil, err
				}
				fmt.Println(result)
			}
			tokens = []Token{}
		}
	}
	return newEnv, nil
}

//take the current environment and write it out to the named file (second parameter)
func store(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameter for STORE")
	}
	switch a2 := t.Right.(type) {
	case *SExpr:
		//should only have a single parameter for STORE
		if a2.Right != NIL {
			return nil, errors.New("shouldn't have more than one parameter for STORE")
		}
		e2, error := evalInner(a2.Left, env)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			f, err := os.Create(string(a3))
			if err != nil {
				return nil, err
			}
			defer f.Close()
			_, err = f.WriteString(TopLevel.String())
			if err != nil {
				return nil, err
			}
			return T, nil
		default:
			return nil, errors.New("STORE parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("Shouldn't get here")
}

func delete(t *SExpr, env Env) (Expr, error) {
	if t.Right == NIL {
		return nil, errors.New("missing parameter for DELETE")
	}
	switch a2 := t.Right.(type) {
	case *SExpr:
		//should only have a single parameter for DELETE
		if a2.Right != NIL {
			return nil, errors.New("shouldn't have more than one parameter for DELETE")
		}
		e2, error := evalInner(a2.Left, env)
		if error != nil {
			return nil, error
		}
		switch a3 := e2.(type) {
		case Atom:
			env.Delete(a3)
			return T, nil
		default:
			return nil, errors.New("DELETE parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("Shouldn't get here")
}

func debug(t *SExpr, env Env) (Expr, error) {
	param, err := nth(1, t)
	if err != nil {
		return nil, err
	}
	v, err := evalInner(param, env)
	if err != nil {
		return nil, err
	}
	global.Log("param is ", param)
	global.Log("v is ", v)
	if v == T {
		global.DEBUG = true
	} else if v == EMPTY {
		global.DEBUG = false
	} else {
		return nil, errors.New("Unknown debug value. Valid values are T and NIL")
	}
	return T, nil
}

func let(t *SExpr, env Env) (Expr, error) {
	//has two params,
	//a list of two-element lists with the scoped variables
	//the command to run with those variables
	//this works like let* does, because why have both?
	//also once a let scope is established, a setq inside of the let will
	//be scoped inside of the let, both for replacing an existing local value
	//or for creating a new one. a setq in a let that refers to a variable in an
	//outer scope will modify that outer scope.
	variables, err := nth(1, t)
	if err != nil {
		return nil, err
	}
	if variables == NIL {
		return nil, errors.New("missing variables for LET")
	}
	l, ok := variables.(*SExpr)
	if !ok {
		return nil, errors.New("LET variable list must be a List")
	}
	innerEnv, err := buildInnerEnv(l, env)
	if err != nil {
		return nil, err
	}
	body, err := nth(2, t)
	if err != nil {
		return nil, err
	}
	return evalInner(body, innerEnv)
}

func buildInnerEnv(l *SExpr, env Env) (Env, error) {
	global.Log("var list == ", l)
	vals := map[Atom]Expr{}
	innerEnv := LocalEnv{Vals: vals, Parent: env}
	i := 0
	for {
		cv, err := nth(i, l)
		if err != nil {
			return nil, err
		}
		if cv == NIL {
			break
		}
		curVar, ok := cv.(*SExpr)
		if !ok {
			return nil, errors.New("LET variable list entry must be a List")
		}
		vn, err := nth(0, curVar)
		if err != nil {
			return nil, err
		}
		varName, ok := vn.(Atom)
		if !ok {
			return nil, errors.New("LET variable names must be Atoms")
		}
		varVal, err := nth(1, curVar)
		if err != nil {
			return nil, err
		}
		varExpr, err := evalInner(varVal, innerEnv)
		if err != nil {
			return nil, err
		}
		innerEnv.Vals[varName] = varExpr
		i++
	}
	return innerEnv, nil
}

//has multiple values, each evaluated one at a time
//returns the last value
func progn(t *SExpr, env Env) (Expr, error) {
	var retval Expr = NIL
	i := 1
	for {
		curParam, err := nth(i, t)
		if err != nil {
			return nil, err
		}
		if curParam == NIL {
			break
		}
		retval, err = evalInner(curParam, env)
		if err != nil {
			return nil, err
		}
		i++
	}
	return retval, nil
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
	le := LocalEnv{Vals: make(map[Atom]Expr), Parent: l.ParentEnv}
	//assign parameter values to parameter names
	count := 0
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
	//call body with new environment
	return evalInner(l.Body, le)
}

//get the nth parameter of the SExpr.
//The function/macro/special form name is the CAR of the SExpr passed in
//pos == 1 for the first parameter. This is the CAR of the CDR of the SExpr passed in
//If there are not enough parameters, NIL is returned
func nth(pos int, e *SExpr) (Expr, error) {
	for i := 0; i < pos; i++ {
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
			}
			//might be numbers, have to do numeric comparison
			r := &big.Rat{}
			r, ok := r.SetString(string(e))
			r2 := &big.Rat{}
			r2, ok2 := r2.SetString(string(e2))
			if ok && ok2 {
				return r.Cmp(r2) == 0
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

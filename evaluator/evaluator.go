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
	"github.com/jonbodner/my_lisp/types"
)

/*
From the paper: "A MICRO-MANUAL FOR LISP - NOT THE WHOLE TRUTH"

found at https://www.cse.sc.edu/~mgv/csce531sp23/micromanualLISP.pdf

DONE 1. value (QUOTE e) = e.
Thus the value of:
(QUOTE A)
is A.

DONE 2. value (CAR e), where value e is a non-empty list, is the first
element of value e.
Thus the value of:
(CAR (QUOTE (A B C)))
is A.

DONE 3. value (CDR e), where value e is a non-empty list, is the list
that remains when the first element of value e is deleted.
Thus the value of:
(CDR (QUOTE (A B C)))
is (B C).

DONE 4. value (CONS e1 e2), is the list that results from prefixing
value e1 onto the list value e2.
Thus the value of:
(CONS (QUOTE A) (QUOTE (B C)))
is (A B C).

DONE 5. value (EQUAL e1 e2) is T if value e1 = value e2. Otherwise, its
value is NIL.
Thus the value of:
(EQUAL (CAR (QUOTE (A B))) (QUOTE A))
is T.

DONE 6. value (ATOM e) is T if value e is an atom; otherwise its value is NIL.

DONE 7. value (COND(p1 e1) ... (pn en)) = value ei, where pi is the first of
the p's whose value is not NIL.
Thus the value of:
(COND
	((ATOM (QUOTE A)) (QUOTE B))
	((QUOTE T) (QUOTE C))
)
is B.

DONE 8. An atom v, regarded as a variable, may have a value.

DONE 9. value ((LAMBDA (v1 ... vn) e) e1 ... en) is the same as value e
but in an environment in which the variables v1 ... vn take the
values of the expressions e1 ... en in the original environment.
Thus the value of:
(
	(LAMBDA (X Y) (CONS (CAR X) Y))
	(QUOTE (A B)) (CDR (QUOTE (C D)))
)
is (A D).

DONE 10. Here's the hard one. value ((LABEL f (LAMBDA (v1 ... vn)
e)) e1 ... en) is the same as value ((LAMBDA (v1 ... vn) e) e1 ... en)
with the additional rule that whenever (f a1 ... an) must be
evaluated, f is replaced by (LABEL f (LAMBDA (v1 ... vn) e)).
Lists beginning with LABEL define functions recursively.
*/

type Evaluator func(*types.SExpr, types.Env) (types.Expr, error)

var TopLevel = make(types.GlobalEnv)

var BuiltIn map[types.Atom]Evaluator

func init() {
	TopLevel[types.T] = types.T
	TopLevel[types.Atom("NIL")] = types.EMPTY

	BuiltIn = map[types.Atom]Evaluator{
		"QUOTE":     quote,
		"CAR":       car,
		"CDR":       cdr,
		"CONS":      cons,
		"ATOM":      atom,
		"EQ":        equal,
		"COND":      cond,
		"LABEL":     label,
		"SETQ":      setq,
		"LAMBDA":    lambda,
		"PROGN":     progn,
		"LET":       let,
		"**DEBUG**": debug,
		"LOAD":      load,
		"STORE":     store,
		"DELETE":    deleteFunc,
	}
}

func Eval(e types.Expr) (types.Expr, error) {
	return evalInner(e, TopLevel)
}

var depth = 0

func evalInner(e types.Expr, env types.Env) (types.Expr, error) {
	depth++
	defer func() {
		global.Log("done depth ", depth)
		depth--
	}()
	global.Log("at depth ", depth)
	global.Log("Evaluating ", e)
	switch t := e.(type) {
	case types.Atom:
		global.Log("\tGot an types.Atom")
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
		return nil, fmt.Errorf("unknown symbol %s ", t)
	case *types.SExpr:
		global.Log("\tGot an types.SExpr")
		switch a := t.Left.(type) {
		case types.Atom:
			global.Log("\t\tLeft is an types.Atom")
			evaluator, ok := BuiltIn[a]
			if ok {
				return evaluator(t, env)
			}
			global.Log("not a builtin")
			//look up variable value in context and process that
			global.Log("looking up ", a)
			expr, ok := env.Get(a)
			if !ok {
				return nil, fmt.Errorf("unknown symbol %s ", a)
			}
			//replace the atom with the value of the expression
			result, err := evalInner(expr, env)
			global.Log("done evaluating")
			if err != nil {
				return nil, err
			}
			newT := &types.SExpr{Left: result, Right: t.Right}
			return evalInner(newT, env)
		case *types.SExpr:
			global.Log("\t\tLeft is an types.SExpr")
			//evaluate the left, then replace left with the evaluated value, and recurse
			lResult, err := evalInner(t.Left, env)
			if err != nil {
				return nil, err
			}
			t.Left = lResult
			return evalInner(t, env)
		case types.Nil:
			global.Log("Got a nil left")
			return t, nil
		case types.Lambda:
			global.Log("\t\tLeft is a types.Lambda")
			return processLambda(a, t, env)
		default:
			return nil, errors.New("shouldn't get here")
		}
	case types.Lambda:
		global.Log("\tGot a lambda")
		return t, nil
	}

	return nil, errors.New("don't know how I got here")
}

func quote(t *types.SExpr, _ types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for QUOTE")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("shouldn't have an types.Atom after a QUOTE")
	case *types.SExpr:
		//should only have a single parameter for QUOTE
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for QUOTE")
		}
		return a2.Left, nil
	default:
		return nil, errors.New("shouldn't get here")
	}
}

func car(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for CAR")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("CAR parameter must be a list")
	case *types.SExpr:
		//should only have a single parameter for CAR
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for CAR")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
			return nil, errors.New("CAR parameter must be a list")
		case *types.SExpr:
			return a3.Left, nil
		}
	default:
		return nil, fmt.Errorf("unknown types.Expr type found: %types.T", a2)
	}
	return nil, errors.New("should never get here")
}

func cdr(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for CDR")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("CDR parameter must be a list")
	case *types.SExpr:
		//should only have a single parameter for CDR
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for CDR")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
			return nil, errors.New("CDR parameter must be a list")
		case *types.SExpr:
			return a3.Right, nil
		}
	default:
		return nil, fmt.Errorf("unknown types.Expr type found: %types.T", a2)
	}
	return nil, errors.New("should never get here")
}

func cons(t *types.SExpr, env types.Env) (types.Expr, error) {
	//must have two params
	//going to construct a types.SExpr out of them
	//first is going to be the left, second is going to be the right
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for CONS")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("CONS parameter must be a list")
	case *types.SExpr:
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		//should have two parameters for CDR
		if a2.Right == types.NIL {
			return nil, errors.New("must have two parameters for CONS")
		}
		switch a3 := a2.Right.(type) {
		case types.Atom:
			return nil, errors.New("CONS parameter must be a list")
		case *types.SExpr:
			if a3.Right != types.NIL {
				return nil, errors.New("must have two parameters for CONS")
			}
			e3, err := evalInner(a3.Left, env)
			if err != nil {
				return nil, err
			}
			if e3 == types.EMPTY {
				e3 = types.NIL
			}
			return &types.SExpr{Left: e2, Right: e3}, nil
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("shouldn't get here")
}

func atom(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for ATOM")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return types.T, nil
	case *types.SExpr:
		//should only have a single parameter for ATOM
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for ATOM")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
			return types.T, nil
		case *types.SExpr:
			if a3.Left == types.NIL && a3.Right == types.NIL {
				return types.T, nil
			}
			return types.EMPTY, nil
		default:
			return nil, errors.New("shouldn't get here")
		}
	}
	return nil, errors.New("shouldn't get here")
}

func equal(t *types.SExpr, env types.Env) (types.Expr, error) {
	//must have two params
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for EQUAL")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("EQUAL parameter must be a list")
	case *types.SExpr:
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		//should have two parameters for EQUAL
		if a2.Right == types.NIL {
			return nil, errors.New("must have two parameters for EQUAL")
		}
		switch a3 := a2.Right.(type) {
		case types.Atom:
			return nil, errors.New("EQUAL parameter must be a list")
		case *types.SExpr:
			if a3.Right != types.NIL {
				return nil, errors.New("must have two parameters for EQUAL")
			}
			e3, err := evalInner(a3.Left, env)
			if err != nil {
				return nil, err
			}
			if isEqual(e2, e3) {
				return types.T, nil
			}
			return types.EMPTY, nil
		}
	default:
		return nil, errors.New("shouldn't get here")
	}
	return nil, errors.New("shouldn't get here")
}

func cond(t *types.SExpr, env types.Env) (types.Expr, error) {
	//find the first non-types.NIL result, and return it
	pos := 1
	for {
		cur, err := nth(pos, t)
		if err != nil {
			return nil, err
		}
		switch cur := cur.(type) {
		case types.Atom:
			return nil, errors.New("cannot have an atom as a COND parameter")
		case types.Nil:
			return types.EMPTY, nil
		case *types.SExpr:
			car, err := evalInner(cur.Left, env)
			if err != nil {
				return nil, err
			}
			if !isEqual(car, types.EMPTY) {
				switch result := cur.Right.(type) {
				case types.Atom:
					return nil, errors.New("cannot have a dotted pair here")
				case types.Nil:
					return types.EMPTY, nil
				case *types.SExpr:
					return evalInner(result.Left, env)
				}
			}
		}
		pos++
	}
}

func label(t *types.SExpr, env types.Env) (types.Expr, error) {
	//must have two params
	//first must be an atom
	//second can be any expression
	//going to assign expression to atom
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for LABEL")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("LABEL parameter must be a list")
	case *types.SExpr:
		//a2.Left must be an atom
		l, ok := a2.Left.(types.Atom)
		if !ok {
			return nil, errors.New("LABEL can only be assigned to an types.Atom")
		}
		//a2.Right must be an *types.SExpr
		a3, ok := a2.Right.(*types.SExpr)
		if !ok {
			return nil, errors.New("LABEL parameter must be a list")
		}
		//a2.Right.Right must be types.NIL
		if a3.Right != types.NIL {
			return nil, errors.New("must have two parameters for LABEL")
		}
		//a2.Right.Left can be anything
		lval := a3.Left
		env.Put(l, lval)
		return l, nil
	}
	return nil, errors.New("shouldn't get here")
}

func setq(t *types.SExpr, env types.Env) (types.Expr, error) {
	//must have two params
	//first must be an atom
	//second can be any expression
	//going to evaluate expression and assign result to atom
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for SETQ")
	}
	switch a2 := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("SETQ parameter must be a list")
	case *types.SExpr:
		//a2.Left must be an atom
		l, ok := a2.Left.(types.Atom)
		if !ok {
			return nil, errors.New("SETQ can only be assigned to an types.Atom")
		}
		//a2.Right must be an *types.SExpr
		a3, ok := a2.Right.(*types.SExpr)
		if !ok {
			return nil, errors.New("SETQ parameter must be a list")
		}
		//a2.Right.Right must be types.NIL
		if a3.Right != types.NIL {
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
	return nil, errors.New("shouldn't get here")
}

func lambda(t *types.SExpr, env types.Env) (types.Expr, error) {
	//must have 2 params
	//param 1 is a list of parameters
	params, err := nth(1, t)
	if err != nil {
		return nil, err
	}
	if params == types.NIL {
		return nil, errors.New("missing parameters for LAMBDA")
	}
	l, ok := params.(*types.SExpr)
	if !ok {
		return nil, errors.New("LAMBDA parameter list must be a List")
	}

	//copy into slice of Atoms
	aList, err := listToSlice(l)
	if err != nil {
		return nil, err
	}

	//param 2 is a types.Expr
	body, err := nth(2, t)
	if err != nil {
		return nil, err
	}
	if body == types.NIL {
		return nil, errors.New("must have two parameters for LAMBDA")
	}
	//a2.Right.Left can be anything
	//returns a new types.Expr type, a types.Lambda, which has its own env
	lambda := types.Lambda{ParentEnv: env, Body: body, Params: aList}
	return lambda, nil
}

// load the environment from the named file. Existing symbols will be overwritten.
func load(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for LOAD")
	}
	switch a2 := t.Right.(type) {
	case *types.SExpr:
		//should only have a single parameter for LOAD
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for LOAD")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
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
			return types.T, nil
		default:
			return nil, errors.New("LOAD parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("shouldn't get here")
}

func internalRepl(r io.Reader) (types.GlobalEnv, error) {
	newEnv := types.GlobalEnv{}
	newEnv[types.T] = types.T
	newEnv[types.Atom("NIL")] = types.EMPTY

	bio := bufio.NewReader(r)
	done := false
	depth := 0
	var tokens []types.Token
	for !done {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			done = true
			continue
		}
		newTokens, newDepth := scanner.Scan(line)
		depth = depth + newDepth
		if depth < 0 {
			return nil, errors.New("invalid -- Too many closing parens")
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
			tokens = []types.Token{}
		}
	}
	return newEnv, nil
}

// take the current environment and write it out to the named file (second parameter)
func store(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for STORE")
	}
	switch a2 := t.Right.(type) {
	case *types.SExpr:
		//should only have a single parameter for STORE
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for STORE")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
			f, err := os.Create(string(a3))
			if err != nil {
				return nil, err
			}
			defer f.Close()
			_, err = f.WriteString(TopLevel.String())
			if err != nil {
				return nil, err
			}
			return types.T, nil
		default:
			return nil, errors.New("STORE parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("shouldn't get here")
}

func deleteFunc(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameter for DELETE")
	}
	switch a2 := t.Right.(type) {
	case *types.SExpr:
		//should only have a single parameter for DELETE
		if a2.Right != types.NIL {
			return nil, errors.New("shouldn't have more than one parameter for DELETE")
		}
		e2, err := evalInner(a2.Left, env)
		if err != nil {
			return nil, err
		}
		switch a3 := e2.(type) {
		case types.Atom:
			env.Delete(a3)
			return types.T, nil
		default:
			return nil, errors.New("DELETE parameter must evaluate to a single value")
		}
	}
	return nil, errors.New("shouldn't get here")
}

func debug(t *types.SExpr, env types.Env) (types.Expr, error) {
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
	if v == types.T {
		global.Debug = true
	} else if v == types.EMPTY {
		global.Debug = false
	} else {
		return nil, errors.New("unknown debug value. Valid values are types.T and types.NIL")
	}
	return types.T, nil
}

func let(t *types.SExpr, env types.Env) (types.Expr, error) {
	//has two params,
	//a list of two-element lists with the scoped variables
	//the command to run with those variables
	//this works like let* does, because why have both?
	//also once a let scope is established, a setq inside the let will
	//be scoped inside the let, both for replacing an existing local value
	//or for creating a new one. a setq in a let that refers to a variable in an
	//outer scope will modify that outer scope.
	variables, err := nth(1, t)
	if err != nil {
		return nil, err
	}
	if variables == types.NIL {
		return nil, errors.New("missing variables for LET")
	}
	l, ok := variables.(*types.SExpr)
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

func buildInnerEnv(l *types.SExpr, env types.Env) (types.Env, error) {
	global.Log("var list == ", l)
	vals := map[types.Atom]types.Expr{}
	innerEnv := types.LocalEnv{Vals: vals, Parent: env}
	i := 0
	for {
		cv, err := nth(i, l)
		if err != nil {
			return nil, err
		}
		if cv == types.NIL {
			break
		}
		curVar, ok := cv.(*types.SExpr)
		if !ok {
			return nil, errors.New("LET variable list entry must be a List")
		}
		vn, err := nth(0, curVar)
		if err != nil {
			return nil, err
		}
		varName, ok := vn.(types.Atom)
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

// has multiple values, each evaluated one at a time
// returns the last value
func progn(t *types.SExpr, env types.Env) (types.Expr, error) {
	var retval types.Expr = types.NIL
	i := 1
	for {
		curParam, err := nth(i, t)
		if err != nil {
			return nil, err
		}
		if curParam == types.NIL {
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

func listToSlice(l *types.SExpr) ([]types.Atom, error) {
	var out []types.Atom

	pos := 0
	for {
		cur, err := nth(pos, l)
		if err != nil {
			return nil, err
		}

		if cur == types.NIL {
			break
		}

		c, ok := cur.(types.Atom)
		if !ok {
			return nil, errors.New("only Atoms can be parameter names")
		}

		out = append(out, c)
		pos++
	}
	return out, nil
}

func processLambda(l types.Lambda, t *types.SExpr, env types.Env) (types.Expr, error) {
	le := types.LocalEnv{Vals: make(map[types.Atom]types.Expr), Parent: l.ParentEnv}
	//assign parameter values to parameter names
	count := 0
	switch paramVals := t.Right.(type) {
	case types.Atom:
		return nil, errors.New("can't have a dotted pair here")
	case types.Nil:
		//do nothing
	case *types.SExpr:
		for count < len(l.Params) {
			var err error
			param, err := nth(count, paramVals)
			if err != nil {
				return nil, err
			}
			if param == types.NIL {
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
		if leftOver != types.NIL {
			return nil, fmt.Errorf("too many parameters for LAMBDA. Expected %d", len(l.Params))
		}
	}
	if count != len(l.Params) {
		return nil, fmt.Errorf("too few parameters for LAMBDA. Expected %d, got %d", len(l.Params), count)
	}
	//call body with new environment
	return evalInner(l.Body, le)
}

// get the nth parameter of the types.SExpr.
// The function/macro/special form name is the CAR of the types.SExpr passed in
// pos == 1 for the first parameter. This is the CAR of the CDR of the types.SExpr passed in
// If there are not enough parameters, types.NIL is returned
func nth(pos int, e *types.SExpr) (types.Expr, error) {
	for i := 0; i < pos; i++ {
		next := e.Right
		switch next := next.(type) {
		case types.Atom:
			return nil, errors.New("can't have a dotted pair here")
		case *types.SExpr:
			e = next
		case types.Nil:
			return types.NIL, nil
		}
	}
	return e.Left, nil
}

func isEqual(e, e2 types.Expr) bool {
	switch e := e.(type) {
	case types.Atom:
		if e2, ok := e2.(types.Atom); ok {
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
	case types.Nil:
		_, ok := e2.(types.Nil)
		return ok
	case *types.SExpr:
		if e2, ok := e2.(*types.SExpr); ok {
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

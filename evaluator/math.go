package evaluator

import (
	"errors"
	"fmt"
	"github.com/jonbodner/my_lisp/global"
	"github.com/jonbodner/my_lisp/types"
	"math/big"
)

func init() {
	BuiltIn[("+")] = plus
	BuiltIn[("-")] = minus
	BuiltIn[("*")] = times
	BuiltIn[("/")] = div
}

func plus(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for + operator")
	}
	params, ok := t.Right.(*types.SExpr)
	if !ok {
		return nil, errors.New("+ parameters must be a list")
	}
	r := &big.Rat{}
	for {
		ev, err := evalInner(params.Left, env)
		if err != nil {
			return nil, err
		}
		global.Log("\tfinished eval -- checking if it's a number")
		r2 := &big.Rat{}
		_, ok = r2.SetString(ev.String())
		if !ok {
			return nil, fmt.Errorf("%s is not a valid number", ev)
		}
		r.Add(r, r2)
		next := params.Right
		if next == types.NIL {
			break
		}
		n, ok := next.(*types.SExpr)
		if !ok {
			return nil, errors.New("+ parameters must be a List")
		}
		params = n
	}
	return types.Atom(r.RatString()), nil
}

func minus(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for - operator")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == types.NIL {
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
		if v == types.NIL {
			//if there was only one value, just negate it
			if first {
				r.Neg(r)
			}
			break
		}
		first = false
	}
	return types.Atom(r.RatString()), nil
}

func times(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for * operator")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == types.NIL {
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
		if v == types.NIL {
			break
		}
	}
	return types.Atom(r.RatString()), nil
}

func div(t *types.SExpr, env types.Env) (types.Expr, error) {
	if t.Right == types.NIL {
		return nil, errors.New("missing parameters for / operator")
	}
	//get first value
	pos := 1
	v, err := nth(pos, t)
	if err != nil {
		return nil, err
	}
	if v == types.NIL {
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
		if v == types.NIL {
			//if there was only one value, just negate it
			if first {
				r.Inv(r)
			}
			break
		}
		first = false
	}
	return types.Atom(r.RatString()), nil
}

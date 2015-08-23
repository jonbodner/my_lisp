package evaluator

import (
    "fmt"
    "errors"
    "math/big"
    . "github.com/jonbodner/my_lisp/types"
    "github.com/jonbodner/my_lisp/global"
)

func init() {
    BuiltIn[Atom("+")] = plus
    BuiltIn[Atom("-")] = minus
    BuiltIn[Atom("*")] = times
    BuiltIn[Atom("/")] = div
}

func plus(t *SExpr, env Env) (Expr, error) {
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
        global.Log("\tfinished eval -- checking if it's a number")
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

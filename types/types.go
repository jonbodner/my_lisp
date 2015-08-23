package types

import (
	"fmt"
	"strings"
	"github.com/jonbodner/my_lisp/global"
)

//expressions

//SExpr is an S-Expression.
//That has a left side that's an expression
type SExpr struct {
	Left  Expr
	Right Expr
}

type Expr interface {
	fmt.Stringer
	isExpr()
}

func (s *SExpr) isExpr() {}

func (s *SExpr) String() string {
	out := "("
	if s.Left != NIL {
		out += s.Left.String()
	}
outer:
	for cur := s.Right; cur != NIL; {
		out += " "
		switch c := cur.(type) {
		case Atom:
			out += ". " + c.String()
			break outer
		case *SExpr:
			out += c.Left.String()
			cur = c.Right
		}
	}
	out += ")"
	return out
}

var EMPTY *SExpr = &SExpr{NIL, NIL}

/*
func (s SExpr) String() string {
	return s.stringInner(true)
}

func (s SExpr) stringInner(top bool) string {
	out := ""
	if top {
		out += "("
	}
	switch l := s.Left.(type) {
	case SExpr:
		out += fmt.Sprintf("%s", l)
	case Atom:
		out += fmt.Sprintf("%s", l)
	}
	switch r := s.Right.(type) {
	case SExpr:
		out += fmt.Sprintf(" %s", r.stringInner(false))
	case Atom:
		out += fmt.Sprintf(" . %s", r)
	}
	if top {
		out = out + ")"
	}
	return out
}
*/

type Atom string

func (a Atom) isExpr() {}
func (a Atom) String() string {
	return string(a)
}

var T Atom = "T"

type Nil struct{}

var NIL Nil

func (n Nil) isExpr() {}
func (n Nil) String() string {
	return "NIL"
}

type Env interface {
	Get(a Atom) (Expr, bool)
	Put(a Atom, e Expr)
}

type GlobalEnv map[Atom]Expr

func (ge GlobalEnv) Get(a Atom) (Expr, bool) {
	global.Log("checking global env for ", a)
	e, ok := ge[a]
	global.Log(e, ok)
	return e, ok
}

func (ge GlobalEnv) Put(a Atom, e Expr) {
	ge[a] = e
}

type LocalEnv struct {
	Vals map[Atom]Expr
	Parent Env
}

func (le LocalEnv) Get(a Atom) (Expr, bool) {
	global.Log("checking local env for ", a)
	e, ok := le.Vals[a]
	if ok {
		global.Log("found ", e, "at my level")
		return e, ok
	}
	global.Log("not in me, going to parent")
	return le.Parent.Get(a)
}

//for now, no shadowing of declarations from outer scopes
//since there's no way to modify the value of a value in an outer scope
//(LABEL and SETQ are both create and assign)
func (le LocalEnv) Put(a Atom, e Expr) {
	//case 1: already defined locally
	if _, ok := le.Vals[a]; ok {
		le.Vals[a] = e
		return
	}
	// case 2: defined somewhere in a parent scope
	if _, ok := le.Parent.Get(a); ok {
		le.Parent.Put(a,e)
		return
	}
	// case 3: never defined
	le.Vals[a] = e
}


type Lambda struct {
	ParentEnv Env
	Params []Atom
	Body Expr
}

func (l Lambda) isExpr() {}
func (l Lambda) String() string {
	sparams := make([]string, len(l.Params))
	for k, v := range l.Params {
		sparams[k] = string(v)
	}
	pstr := strings.Join(sparams," ")

	return "(LAMBDA (" + pstr +") "+l.Body.String()+" )"
}


//tokens

type Token interface {
	fmt.Stringer
	TokenForm() string
}

type LParen struct{}

var LPAREN LParen

func (l LParen) TokenForm() string { return "(" }
func (l LParen) String() string {
	return "LPAREN"
}

type RParen struct{}

var RPAREN RParen

func (r RParen) TokenForm() string { return ")" }
func (r RParen) String() string {
	return "RPAREN"
}

type Dot struct{}

var DOT Dot

func (d Dot) TokenForm() string { return "." }
func (d Dot) String() string {
	return "DOT"
}

type Quote struct{}

var QUOTE Quote

func (q Quote) TokenForm() string { return "'" }
func (q Quote) String() string {
	return "QUOTE"
}

type NAME string

func (n NAME) String() string    { return string(n) }
func (n NAME) TokenForm() string { return string(n) }

package main

import (
	"fmt"
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

//tokens

type Token interface {
	fmt.Stringer
	tokenForm() string
}

type LParen struct{}

var LPAREN LParen

func (l LParen) tokenForm() string { return "(" }
func (l LParen) String() string {
	return "LPAREN"
}

type RParen struct{}

var RPAREN RParen

func (r RParen) tokenForm() string { return ")" }
func (r RParen) String() string {
	return "RPAREN"
}

type Dot struct{}

var DOT Dot

func (d Dot) tokenForm() string { return "." }
func (d Dot) String() string {
	return "DOT"
}

type Quote struct{}

var QUOTE Quote

func (q Quote) tokenForm() string { return "'" }
func (q Quote) String() string {
	return "QUOTE"
}

type NAME string

func (n NAME) String() string    { return string(n) }
func (n NAME) tokenForm() string { return string(n) }

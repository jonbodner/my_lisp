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
	isExpr()
}

func (s SExpr) isExpr() {}
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

type Atom string

func (a Atom) isExpr() {}
func (a Atom) String() string {
	return string(a)
}

type Nil struct{}

var NIL Nil

func (n Nil) isExpr() {}
func (n Nil) String() string {
	return "NIL"
}

//tokens

type Token interface {
	isTok()
}

type LParen struct{}

var LPAREN LParen

func (l LParen) isTok() {}
func (l LParen) String() string {
	return "LPAREN"
}

type RParen struct{}

var RPAREN RParen

func (r RParen) isTok() {}
func (r RParen) String() string {
	return "RPAREN"
}

type Dot struct{}

var DOT Dot

func (d Dot) isTok() {}
func (d Dot) String() string {
	return "DOT"
}

type Quote struct{}

var QUOTE Quote

func (q Quote) isTok() {}
func (q Quote) String() string {
	return "QUOTE"
}

type NAME string

func (n NAME) isTok() {}

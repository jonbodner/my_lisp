package main

import (
	"fmt"
	"testing"
)

func TestSExpr(t *testing.T) {
	emptyList := SExpr{Left: NIL, Right: NIL}
	fmt.Println(emptyList)
	if emptyList.String() != "()" {
		t.Fail()
	}

	singleItemList := SExpr{Left: Atom("WORD"), Right: NIL}
	fmt.Println(singleItemList)
	if singleItemList.String() != "(WORD)" {
		t.Fail()
	}

	dottedList := SExpr{Left: Atom("WORD"), Right: Atom("WORD2")}
	fmt.Println(dottedList)
	if dottedList.String() != "(WORD . WORD2)" {
		t.Fail()
	}

	twoItemList := SExpr{Left: Atom("WORD"), Right: SExpr{Left: Atom("WORD2"), Right: NIL}}
	fmt.Println(twoItemList)
	if twoItemList.String() != "(WORD WORD2)" {
		t.Fail()
	}

	nestedListStart := SExpr{Left: SExpr{Left: Atom("WORD"), Right: NIL}, Right: NIL}
	fmt.Println(nestedListStart)
	if nestedListStart.String() != "((WORD))" {
		t.Fail()
	}

	nestedListStartEnd := SExpr{Left: SExpr{Left: Atom("WORD"), Right: SExpr{Left: Atom("WORD2"), Right: NIL}}, Right: SExpr{Left: Atom("WORD3"), Right: SExpr{Left: Atom("WORD4"), Right: NIL}}}
	fmt.Println(nestedListStartEnd)
	if nestedListStartEnd.String() != "((WORD WORD2) WORD3 WORD4)" {
		t.Fail()
	}
}

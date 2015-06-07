package evaluator

import (
    "testing"
    "host.bodnerfamily.com/my_lisp/scanner"
    "host.bodnerfamily.com/my_lisp/parser"
)

//additional
func TestMinus(t *testing.T) {
    internalEvaluator(t, "( - 1)", "-1")
    internalEvaluator(t, "( - 1 1)", "0")
    internalEvaluator(t, "( - 1 5)", "-4")
    internalEvaluator(t, "( - 5 1)", "4")
    internalEvaluator(t, "( - 1 1 1 1 1 1)", "-4")
}

func TestPlus(t *testing.T) {
    internalEvaluator(t, "( + 1)", "1")
    internalEvaluator(t, "( + 1 1)", "2")
    internalEvaluator(t, "( + 1 5)", "6")
    internalEvaluator(t, "( + 5 1)", "6")
    internalEvaluator(t, "( + 1 1 1 1 1 1)", "6")
}

func TestTimes(t *testing.T) {
    internalEvaluator(t, "( + 1)", "1")
    internalEvaluator(t, "( + 1 1)", "2")
    internalEvaluator(t, "( + 1 5)", "6")
    internalEvaluator(t, "( + 5 1)", "6")
    internalEvaluator(t, "( + 1 1 1 1 1 1)", "6")
}

func TestDivide(t *testing.T) {
    internalEvaluator(t, "( + 1)", "1")
    internalEvaluator(t, "( + 1 1)", "2")
    internalEvaluator(t, "( + 1 5)", "6")
    internalEvaluator(t, "( + 5 1)", "6")
    internalEvaluator(t, "( + 1 1 1 1 1 1)", "6")
}

//core
func TestQuote(t *testing.T) {

}

func TestCar(t *testing.T) {

}

func TestCdr(t *testing.T) {

}

func TestCons(t *testing.T) {

}

func TestEq(t *testing.T) {

}

func TestAtom(t *testing.T) {

}

func TestCond(t *testing.T) {

}

func TestEnv(t *testing.T) {

}

func TestLambda(t *testing.T) {

}

func TestLabel(t *testing.T) {

}

func TestSetq(t *testing.T) {

}

func internalEvaluator(t *testing.T, input string, expected string) {
    tokens, _ := scanner.Scan(input)
    expr, _, _ := parser.Parse(tokens)
    out, err := Eval(expr)
    if err != nil {
        t.Errorf("Unexpected error %v", err)
    } else {
        if out.String() != expected {
            t.Errorf("Expected %s, got %v ", expected, out)
        }
    }
}
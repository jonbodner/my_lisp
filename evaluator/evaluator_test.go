package evaluator

import (
	"github.com/jonbodner/my_lisp/parser"
	"github.com/jonbodner/my_lisp/scanner"
	"testing"
)

// additional
func TestMinus(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{"unary", "( - 1)", "-1"},
		{"zero", "( - 1 1)", "0"},
		{"negative", "( - 1 5)", "-4"},
		{"positive", "(- 5 1)", "4"},
		{"repeated", "( - 1 1 1 1 1 1)", "-4"},
		{"error", "(-)", "missing parameters for - operator"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			internalEvaluator(t, d.input, d.expected)
		})
	}
}

func TestPlus(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{"unary", "( + 1)", "1"},
		{"positive", "( + 1 1)", "2"},
		{"positive2", "( + 1 5)", "6"},
		{"positive3", "( + 5 1)", "6"},
		{"repeated", "( + 1 1 1 1 1 1)", "6"},
		{"error", "(+)", "missing parameters for + operator"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			internalEvaluator(t, d.input, d.expected)
		})
	}
}

func TestTimes(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{"unary", "( * 1)", "1"},
		{"times", "( * 1 1)", "1"},
		{"times2", "( * 1 5)", "5"},
		{"repeated", "( * 1 1 1 1 1 1)", "1"},
		{"error", "(*)", "missing parameters for * operator"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			internalEvaluator(t, d.input, d.expected)
		})
	}
}

func TestDivide(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{"unary", "( / 1)", "1"},
		{"divide", "( / 1 1)", "1"},
		{"divide2", "( / 1 5)", "1/5"},
		{"divide3", "( / 5 1)", "5"},
		{"repeated", "( / 1 1 1 1 1 1)", "1"},
		{"error", "(/)", "missing parameters for / operator"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			internalEvaluator(t, d.input, d.expected)
		})
	}
}

// core
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
		if err.Error() != expected {
			t.Errorf("Unexpected error. Expected %s, got %v", expected, err.Error())
		}
	} else {
		if out.String() != expected {
			t.Errorf("Expected %s, got %v ", expected, out)
		}
	}
}

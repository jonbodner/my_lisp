package parser

import "github.com/jonbodner/my_lisp/types"

type ParseError struct {
	msg    string
	tokens []types.Token
	pos    int
}

func (te ParseError) Error() string {
	out := te.msg
	tokenPart := te.buildTokenForm()
	if len(tokenPart) > 0 {
		out += ": "
		out += tokenPart
	}
	return out
}

func (te ParseError) buildTokenForm() string {
	out := ""
	for i := 0; i < len(te.tokens); i++ {
		if i == te.pos {
			out += "_"
		}
		out += te.tokens[i].TokenForm()
		if i == te.pos {
			out += "_"
		}
		out += " "
	}
	return out
}

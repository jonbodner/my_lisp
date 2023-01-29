package scanner

import "github.com/jonbodner/my_lisp/types"

func Scan(s string) ([]types.Token, int) {
	var out []types.Token
	var curTokenTxt []rune
	buildCurToken := func() {
		if len(curTokenTxt) > 0 {
			if len(curTokenTxt) == 1 && curTokenTxt[0] == '.' {
				out = append(out, types.DOT)
			} else {
				out = append(out, types.NAME(curTokenTxt))
			}
			curTokenTxt = make([]rune, 0)
		}
	}
	update := func(t types.Token) {
		buildCurToken()
		out = append(out, t)
	}

	depth := 0
	for _, c := range s {
		switch c {
		case '(':
			update(types.LPAREN)
			depth++
		case ')':
			update(types.RPAREN)
			depth--
		case '.':
			curTokenTxt = append(curTokenTxt, c)
		case '\n', '\r', '\t', ' ':
			buildCurToken()
		case '\'':
			update(types.QUOTE)
		default:
			curTokenTxt = append(curTokenTxt, c)
		}
	}
	buildCurToken()
	return out, depth
}

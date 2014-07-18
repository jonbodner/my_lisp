package main

func Scan(s string) ([]Token, int) {
	out := make([]Token, 0)
	curTokenTxt := make([]rune, 0)
	buildCurToken := func() {
		if len(curTokenTxt) > 0 {
			out = append(out, NAME(string(curTokenTxt)))
			curTokenTxt = make([]rune, 0)
		}
	}
	update := func(t Token) {
		buildCurToken()
		out = append(out, t)
	}

	depth := 0
	for _, c := range s {
		switch c {
		case '(':
			update(LPAREN)
			depth++
		case ')':
			update(RPAREN)
			depth--
		case '.':
			update(DOT)
		case '\n', '\r', '\t', ' ':
			buildCurToken()
		case '\'':
			update(QUOTE)
		default:
			curTokenTxt = append(curTokenTxt, c)
		}
	}
	buildCurToken()
	return out, depth
}

package reggi

type symbol int

const (
	AnyChar symbol = iota
	Pipe
	LParen
	RParen
	Char
	ZeroOrMore
	OneOrMore
	ZeroOrOne
)

type token struct {
	symbol symbol
	char   rune
}

// description
func lex(input string) []token {
	var tokens []token
	for _, r := range input {
		tokens = append(tokens, lexRune(r))
	}

	return tokens
}

func lexRune(r rune) token {
	var t token
	switch r {
	case '(':
		t.symbol = LParen
	case ')':
		t.symbol = RParen
	case '|':
		t.symbol = Pipe
	case '.':
		t.symbol = AnyChar
	case '*':
		t.symbol = ZeroOrMore
	case '?':
		t.symbol = ZeroOrOne
	case '+':
		t.symbol = OneOrMore
	default:
		t.symbol = Char
		t.char = r
	}

	return t
}

package reggi

type Parser struct {
	tokens []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p Parser) Parse() Node {
	root := Group{}

	for _, t := range p.tokens {
		switch t.symbol {
		case Char:
			root.Append(&CharLiteral{Char: t.char})
		case AnyChar:
			root.Append(&Wildcard{})
		}
	}

	return &root
}

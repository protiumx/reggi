package reggi

type Parser struct {
	tokens []token
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p Parser) Parse() Node {
	var root CompositeNode
	root = &Group{}

	for _, t := range p.tokens {
		switch t.symbol {
		case Char:
			root.Append(&CharLiteral{Char: t.char})
		case AnyChar:
			root.Append(&Wildcard{})
		case Pipe:
			switch r := root.(type) {
			case *Branch:
				r.Split()
			default:
				root = &Branch{Nodes: []Node{root, &Group{}}}
			}
		}
	}

	return root
}

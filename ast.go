package reggi

type Node interface {
	compile() (head *State, tail *State)
}

type CompositeNode interface {
	Node
	Append(node Node)
}

type Group struct {
	Nodes []Node
}

func (g *Group) Append(node Node) {
	g.Nodes = append(g.Nodes, node)
}

func (g *Group) compile() (*State, *State) {
	head := &State{}
	tail := head

	for _, n := range g.Nodes {
		nextHead, nextTail := n.compile()
		tail.merge(nextHead)
		tail = nextTail
	}

	return head, tail
}

type CharLiteral struct {
	Char rune
}

func (c *CharLiteral) compile() (*State, *State) {
	head := &State{}
	tail := &State{}
	head.addTransition(tail, Predicate{allowed: string(c.Char)}, string(c.Char))

	return head, tail
}

type Wildcard struct{}

func (*Wildcard) compile() (*State, *State) {
	head := &State{}
	tail := &State{}
	head.addTransition(tail, Predicate{disallowed: "\n"}, ".")

	return head, tail
}

type Branch struct {
	Nodes []Node
}

func (b *Branch) Append(node Node) {
}

func (b *Branch) compile() (head *State, tails State) {
	panic("")
}

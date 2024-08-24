package reggi

import (
	"fmt"
	"strings"
)

type Node interface {
	compile() (head *State, tail *State)
	string(indent int) string
}

type CompositeNode interface {
	Node
	Append(node Node)
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

func (c *CharLiteral) string(indent int) string {
	padding := strings.Repeat("--", indent)
	return fmt.Sprintf("%sCharLiteral('%c')", padding, c.Char)
}

type Wildcard struct{}

func (*Wildcard) compile() (*State, *State) {
	head := &State{}
	tail := &State{}
	head.addTransition(tail, Predicate{disallowed: "\n"}, ".")

	return head, tail
}

func (w *Wildcard) string(indent int) string {
	return strings.Repeat("--", indent) + "Wildcard"
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

func (g *Group) string(indent int) string {
	return compileToString("Group", g.Nodes, indent)
}

func (g *Group) String() string {
	return "\n" + g.string(1)
}

type Branch struct {
	Nodes []Node
}

func (b *Branch) Append(node Node) {
	// Append to the last CompositeNode
	for i := len(b.Nodes) - 1; i >= 0; i-- {
		switch n := b.Nodes[i].(type) {
		case CompositeNode:
			n.Append(node)
			return
		}
	}

	panic("branch should have at least one composite node")
}

func (b *Branch) compile() (head *State, tails *State) {
	head = &State{}

	for _, node := range b.Nodes {
		next, _ := node.compile()
		head.merge(next)
	}
	return head, head
}

// Split adds a new branch
func (b *Branch) Split() {
	b.Nodes = append(b.Nodes, &Group{})
}

func (b *Branch) string(indent int) string {
	return compileToString("Branch", b.Nodes, indent)
}

func (b *Branch) String() string {
	return "\n" + b.string(1)
}

func compileToString(s string, nodes []Node, indent int) string {
	padding := strings.Repeat("--", indent)
	res := padding + s
	for _, node := range nodes {
		res += fmt.Sprintf("\n%s%s", padding, node.string(indent+1))
	}

	return res
}

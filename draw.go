package reggi

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var NormalStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("133"))

var SuccessStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("71"))

var FailStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("9"))

var SymbolStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("111"))

func visitNodes(node *State, transitions *OrderedSet[Transition], visited *OrderedSet[*State]) {
	if visited.has(node) {
		return
	}

	transitions.add(node.transitions...)
	visited.add(node)

	for _, t := range node.transitions {
		visitNodes(t.to, transitions, visited)
	}
}

func StatusStyle(status Status, s string) string {
	switch status {
	case StatusNormal:
		return NormalStyle.Render(s)
	case StatusSuccess:
		return SuccessStyle.Render(s)
	case StatusFail:
		return FailStyle.Render(s)
	default:
		return s
	}
}

func (r *Runner) snapshot() string {
	s := r.root
	transitions := OrderedSet[Transition]{}
	nodes := OrderedSet[*State]{}

	visitNodes(s, &transitions, &nodes)

	status := r.status()
	out := make([]string, 0, len(nodes.set))

	for i, t := range transitions.list() {
		from, to := nodes.index(t.from), nodes.index(t.to)

		fromStr := "(" + strconv.Itoa(from) + ")"
		toStr := "(" + strconv.Itoa(to) + ")"

		// account for initial or current state
		if (i == 0 && status == StatusFail) || t.from == r.current {
			fromStr = StatusStyle(status, fromStr)
		} else if t.to == r.current {
			toStr = StatusStyle(status, toStr)
		}

		out = append(out, fromStr+"--"+SymbolStyle.Render(t.debugSym))
		if i == len(transitions.set)-1 {
			out = append(out, toStr)
		}
	}

	return strings.Join(out, "--")
}

package reggi

import (
	"fmt"
	"sort"
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
	Underline(true).
	Bold(true).
	Foreground(lipgloss.Color("111"))

// depth-first traversal
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
	return strings.Join(r.buildLines(), "\n")
}

// buildLines traverses the states in depth-first order and builds the output line by line
// accounting for padding and final states
func (r *Runner) buildLines() []string {
	lines := make([]string, 1)
	stack := []*State{r.root}
	nodeID := 0
	currentLine := 0

	status := r.status()

	for len(stack) > 0 {
		last := len(stack) - 1
		node := stack[last]
		stack = stack[:last]

		current := fmt.Sprintf("(%d)", nodeID)
		// color if active
		if r.activeStates.has(node) {
			current = StatusStyle(status, current)
		}

		// handle node connection (from)--(to)
		if lines[currentLine] != "" {
			lines[currentLine] += "--"
		}
		lines[currentLine] += current
		nodeID++

		// handle leaf node, i.e. node without transitions
		if len(node.transitions) == 0 {
			// we will process the next level
			currentLine++
			continue
		}

		// cater space for all levels
		for len(lines) < len(node.transitions) {
			lines = append(lines, "")
		}

		for i := len(node.transitions) - 1; i >= 0; i-- {
			// draw in the correct level
			c := currentLine + i
			t := node.transitions[i]

			// connect to parent
			if i > 0 && lines[c] == "" {
				lines[i] += "  \\"
			}

			lines[c] += "--" + t.debugSym
			stack = append(stack, t.to)
		}
	}

	return lines
}

func (r *Runner) activeSymbols() string {
	if r.status() != StatusNormal {
		return ""
	}

	symbols := ""
	for i, s := range r.activeStates.list() {
		if i > 0 {
			symbols += ", "
		}
		symbols += s.transitions[0].debugSym
	}

	return symbols
}

func sortVisitedStates(states []*State, nodes OrderedSet[*State]) []*State {
	sort.Slice(states, func(i, j int) bool {
		return nodes.index(states[i]) < nodes.index(states[j])
	})

	return states
}

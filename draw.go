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

		// handle leaf node, i.e. node without transitions
		if len(node.transitions) == 0 {
			leaf := fmt.Sprintf("(%d)", nodeID)
			if r.activeStates.has(node) {
				leaf = StatusStyle(status, leaf)
			}
			lines[currentLine] += "--" + leaf

			// we don't want an extra empty line if there no more nodes
			if len(stack) > 0 {
				currentLine++
				lines = append(lines, strings.Repeat(" ", 5)+"\\")
			}

			nodeID++
			continue
		}

		// handle node connection
		if lines[currentLine] != "" {
			lines[currentLine] += "--"
		}

		current := fmt.Sprintf("(%d)", nodeID)

		// handle node coloring
		if r.activeStates.has(node) {
			current = StatusStyle(status, current)
		}

		lines[currentLine] += current + "--" + node.transitions[0].debugSym

		for i := len(node.transitions) - 1; i >= 0; i-- {
			stack = append(stack, node.transitions[i].to)
		}

		nodeID++
	}

	return lines
}

func (r *Runner) info() string {
	status := r.status()
	if status == StatusFail {
		return FailStyle.Render("Failed")
	}

	if status == StatusSuccess {
		return SuccessStyle.Render("Matched")
	}

	symbols := ""
	for i, s := range r.activeStates.list() {
		if i > 0 {
			symbols += ", "
		}
		symbols += s.transitions[0].debugSym
	}

	return fmt.Sprintf("Trying %s", SuccessStyle.Render(symbols))
}

func sortVisitedStates(states []*State, nodes OrderedSet[*State]) []*State {
	sort.Slice(states, func(i, j int) bool {
		return nodes.index(states[i]) < nodes.index(states[j])
	})

	return states
}

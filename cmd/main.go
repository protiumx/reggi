package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	reggi "protiumx.dev/reggi"
)

var help = lipgloss.NewStyle().
	Foreground(lipgloss.Color("248")).
	Render("h - prev transition\nl - next transition")

var (
	resultMatch   = reggi.SuccessStyle.Render("match")
	resultNoMatch = reggi.SuccessStyle.Render("no match")
)

type model struct {
	step  int
	steps []reggi.DebugStep
	input string
	regex string
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	step := m.steps[m.step]
	text := strings.Split(m.input, "")
	last := m.step == len(m.steps)-1

	if last && step.Status != reggi.StatusSuccess {
		for i := range text {
			text[i] = reggi.FailStyle.Render(text[i])
		}

		return m.render(text, step, last)
	}

	for i := 0; i < step.Offset; i++ {
		text[i] = reggi.FailStyle.Render(text[i])
	}

	for i := step.Offset; i < step.CurrentIndex; i++ {
		text[i] = reggi.SuccessStyle.Render(text[i])
	}

	if step.CurrentIndex < len(text) {
		text[step.CurrentIndex] = reggi.SymbolStyle.Render(text[step.CurrentIndex])
	}

	if last {
		text[len(text)-1] = reggi.StatusStyle(step.Status, text[len(text)-1])
	}

	return m.render(text, step, last)
}

func (m *model) render(text []string, step reggi.DebugStep, last bool) string {
	inputInfo := fmt.Sprintf("input: %q\tregex: %q", m.input, m.regex)
	stepInfo := fmt.Sprintf("\tstep: %d/%d\n\n", m.step+1, len(m.steps))
	stateInfo := strings.Join(text, "  ") + "\n" + step.Snapshot

	if last {
		result := resultMatch
		if step.Status != reggi.StatusSuccess {
			result = resultNoMatch
		}
		stateInfo += " - " + result
	}

	return inputInfo + stepInfo + stateInfo + "\n\n" + help
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h", "left":
			if m.step > 0 {
				m.step--
			}

		// The "down" and "j" keys move the cursor down
		case "l", "right":
			if m.step < len(m.steps)-1 {
				m.step++
			}
		}
	}

	return m, nil
}

func main() {
	if len(os.Args) < 3 {
		panic("usage: reggi <regex> <input>")
	}

	reggex, input := os.Args[1], os.Args[2]
	reggi := reggi.NewReggi(reggex)
	steps := reggi.DebugMatch(input)

	m := model{
		steps: steps,
		input: input,
		regex: reggex,
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}

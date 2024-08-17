package reggi

import "testing"

func abcBuilder() *State {
	state1, state2, state3, state4 := &State{}, &State{}, &State{}, &State{}

	state1.addTransition(state2, Predicate{allowed: "a"}, "a")
	state2.addTransition(state3, Predicate{allowed: "b"}, "b")
	state3.addTransition(state4, Predicate{allowed: "c"}, "c")
	return state1
}

func Test_Snapshot(t *testing.T) {
	type test struct {
		name       string
		input      string
		expected   string
		fsmBuilder func() *State
	}

	// NOTE: lipgloss uses https://github.com/muesli/termenv to determine the TTY profile.
	// When running tests, i.e. not in a terminal, the strings are rendered as normal ASCII
	tests := []test{
		{
			name:       "initial",
			input:      "",
			fsmBuilder: abcBuilder,
			expected:   `(0)--a--(1)--b--(2)--c--(3)`,
		},
		{
			name:       "after single char",
			input:      "a",
			fsmBuilder: abcBuilder,
			expected:   `(0)--a--(1)--b--(2)--c--(3)`,
		},
		{
			name:       "full match",
			input:      "abc",
			fsmBuilder: abcBuilder,
			expected:   `(0)--a--(1)--b--(2)--c--(3)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRunner(tt.fsmBuilder())
			for _, char := range tt.input {
				runner.Next(char)
			}

			if s := runner.snapshot(); s != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, s)
			}
		})
	}
}

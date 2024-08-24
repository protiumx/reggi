package reggi

import "testing"

func regexFSM(regex string) *State {
	return NewReggi(regex).fms
}

func Test_Snapshot(t *testing.T) {
	type test struct {
		name      string
		input     string
		expected  string
		rootState *State
	}

	// NOTE: lipgloss uses https://github.com/muesli/termenv to determine the TTY profile.
	// When running tests, i.e. not in a terminal, the strings are rendered as normal ASCII
	tests := []test{
		{
			name:      "simple match",
			input:     "abc",
			rootState: regexFSM("abc"),
			expected:  `(0)--a--(1)--b--(2)--c--(3)`,
		},
		{
			name:      "branches",
			input:     "abc",
			rootState: regexFSM("ab|ac"),
			expected:  "(0)--a--(1)--b--(2)\n  \\--a--(3)--c--(4)",
		},
		{
			name:      "branches",
			input:     "aaa",
			rootState: regexFSM("a|b|c"),
			expected:  "(0)--a--(1)\n  \\--b--(2)\n  \\--c--(3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRunner(tt.rootState)
			for _, char := range tt.input {
				runner.Next(char)
				runner.Start()
			}

			if s := runner.snapshot(); s != tt.expected {
				t.Fatalf("Expected drawing to be \n\"%s\", got\n\"%s\"", tt.expected, s)
			}
		})
	}
}

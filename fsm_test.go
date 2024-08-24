package reggi

import (
	"regexp"
	"strings"
	"testing"
)

func TestMatchString(t *testing.T) {
	testCases := []struct {
		name, regex, input string
	}{
		{
			name:  "empty input",
			regex: "abc",
			input: "",
		},
		{
			name:  "empty regex",
			regex: "",
			input: "abc",
		},
		{
			name:  "no match",
			regex: "xxx",
			input: "abc",
		},
		{
			name:  "partial match",
			regex: "ab",
			input: "abc",
		},
		{
			name:  "match",
			regex: "abc",
			input: "abc",
		},
		{
			name:  "nested",
			regex: "a(b(d))c",
			input: "abdc",
		},
		{
			name:  "substring match with reset",
			regex: "aA",
			input: "aaA",
		},
		{
			name:  "substring match without reset",
			regex: "B",
			input: "ABA",
		},
		{
			name:  "multi byte",
			regex: "Ȥ",
			input: "Ȥ",
		},
		{
			name:  "complex multibyte characters",
			regex: string([]byte{0xcc, 0x87, 0x30}),
			input: string([]byte{0xef, 0xbf, 0xbd, 0x30}),
		},
		// wildcards
		{
			name:  "wildcard matching",
			regex: "ab.",
			input: "abc",
		},
		{
			name:  "wildcard not matching",
			regex: "ab.",
			input: "ab",
		},
		{
			name:  "wildcard with new lines",
			regex: "..0",
			input: "0\n0",
		},
		{
			name:  "branch matches first",
			regex: "ab|cd",
			input: "ab",
		},
		{
			name:  "branch matches last",
			regex: "ab|cd",
			input: "cd",
		},
		{
			name:  "branch no match",
			regex: "ab|cd",
			input: "xx",
		},
		{
			name:  "branch with shared chars",
			regex: "dog|dot",
			input: "dog",
		},
		{
			name:  "branch with shared chars",
			regex: "dog|dot",
			input: "dot",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			compareWithStd(t, NewReggi(tt.regex), tt.input, tt.regex)
		})
	}
}

func compareWithStd(t *testing.T, r *Reggi, input, regex string) {
	t.Helper()

	result := r.MatchString(input)
	goResult := regexp.MustCompile(regex).MatchString(input)

	if result != goResult {
		t.Errorf(
			"Mismatch - Regex=%s (bytes=%x) Input=%s (byte=%x) Result=%v GoResult=%v",
			regex,
			[]byte(regex),
			input,
			[]byte(input),
			result,
			goResult,
		)
	}
}

func FuzzFSM(f *testing.F) {
	f.Add("abc", "abc")
	f.Add("abc", "")
	f.Add("abc", "xxx")
	f.Add("ca(t)(s)", "dog")

	f.Fuzz(func(t *testing.T, regex, input string) {
		if strings.ContainsAny(regex, "[]{}$^|*+?\\") {
			t.Skip()
		}

		_, err := regexp.Compile(regex)
		if err != nil {
			t.Skip()
		}

		compareWithStd(t, NewReggi(regex), input, regex)
	})
}

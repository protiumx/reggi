package reggi

import (
	"regexp"
	"strings"
	"testing"
)

func TestMatchString(t *testing.T) {
	testCases := []struct {
		name, input, regex string
	}{
		{
			name:  "empty input",
			input: "",
			regex: "abc",
		},
		{
			name:  "empty regex",
			input: "abc",
			regex: "",
		},
		{
			name:  "no match",
			input: "abc",
			regex: "xxx",
		},
		{
			name:  "partial match",
			input: "abc",
			regex: "ab",
		},
		{
			name:  "match",
			input: "abc",
			regex: "abc",
		},
		{
			name:  "nested",
			input: "abdc",
			regex: "a(b(d))c",
		},
		{
			name:  "substring match with reset",
			input: "aaA",
			regex: "aA",
		},
		{
			name:  "substring match without reset",
			input: "ABA",
			regex: "B",
		},
		{
			name:  "multi byte",
			input: "Ȥ",
			regex: "Ȥ",
		},
		{
			name:  "complex multibyte characters",
			input: string([]byte{0xef, 0xbf, 0xbd, 0x30}),
			regex: string([]byte{0xcc, 0x87, 0x30}),
		},
		// wildcards
		{
			name:  "wildcard matching",
			input: "abc",
			regex: "ab.",
		},
		{
			name:  "wildcard not matching",
			input: "ab",
			regex: "ab.",
		},
		{
			name:  "wildcard with new lines",
			input: "0\n0",
			regex: "..0",
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

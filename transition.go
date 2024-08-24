package reggi

import "strings"

type Predicate struct {
	allowed    string
	disallowed string
}

func (p *Predicate) test(in rune) bool {
	if p.allowed != "" && p.disallowed != "" {
		panic("must be mutual exclusive")
	}

	if p.allowed != "" {
		return strings.ContainsRune(p.allowed, in)
	}

	if p.disallowed != "" {
		return !strings.ContainsRune(p.disallowed, in)
	}

	return false
}

type Transition struct {
	debugSym string

	from      *State
	to        *State
	predicate Predicate
}

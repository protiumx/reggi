package reggi

type Status string

const (
	StatusNormal  Status = "normal"
	StatusSuccess Status = "success"
	StatusFail    Status = "fail"
)

type State struct {
	transitions []Transition
}

func (s *State) firstMatchingTransition(input rune) *State {
	for _, t := range s.transitions {
		if t.predicate.test(input) {
			return t.to
		}
	}

	return nil
}

func (s *State) addTransition(to *State, predicate Predicate, debugSym string) {
	s.transitions = append(s.transitions, Transition{
		debugSym:  debugSym,
		to:        to,
		from:      s,
		predicate: predicate,
	})
}

func (s *State) isSuccess() bool {
	return len(s.transitions) == 0
}

func (s *State) merge(s2 *State) {
	for _, t := range s2.transitions {
		s.addTransition(t.to, t.predicate, t.debugSym)
	}

	s2.delete()
}

func (s *State) delete() {
	s.transitions = nil
}

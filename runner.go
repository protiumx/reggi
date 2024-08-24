package reggi

type Runner struct {
	root         *State
	activeStates Set[*State]
}

func NewRunner(root *State) *Runner {
	return &Runner{root: root, activeStates: NewSet(root)}
}

func (r *Runner) Next(input rune) {
	if r.activeStates.size() == 0 {
		return
	}

	nextActiveStates := Set[*State]{}
	for state := range r.activeStates {
		for _, nextState := range state.matchingTransitions(input) {
			nextActiveStates.add(nextState)
		}
	}

	r.activeStates = nextActiveStates
}

func (r *Runner) Reset() {
	r.activeStates = NewSet(r.root)
}

func (r *Runner) status() Status {
	if r.activeStates.size() == 0 {
		return StatusFail
	}

	for state := range r.activeStates {
		if state.isSuccess() {
			return StatusSuccess
		}
	}

	return StatusNormal
}

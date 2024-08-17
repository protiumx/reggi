package reggi

type Runner struct {
	root    *State
	current *State
}

func NewRunner(root *State) *Runner {
	return &Runner{root: root, current: root}
}

func (r *Runner) Next(input rune) {
	if r.current == nil {
		return
	}

	r.current = r.current.firstMatchingTransition(input)
}

func (r *Runner) Reset() {
	r.current = r.root
}

func (r *Runner) status() Status {
	if r.current == nil {
		return StatusFail
	}

	if r.current.isSuccess() {
		return StatusSuccess
	}

	return StatusNormal
}

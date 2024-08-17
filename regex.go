package reggi

type Reggi struct {
	raw string
	fms *State
}

type DebugStep struct {
	Snapshot     string
	Status       Status
	CurrentIndex int
	Offset       int
}

func NewReggi(regex string) *Reggi {
	tokens := lex(regex)
	parser := NewParser(tokens)
	node := parser.Parse()
	state, _ := node.compile()
	return &Reggi{fms: state, raw: regex}
}

func (r *Reggi) MatchString(input string) bool {
	runner := NewRunner(r.fms)
	return match(runner, []rune(input), nil, 0)
}

func (r *Reggi) DebugMatch(input string) []DebugStep {
	runner := NewRunner(r.fms)
	debug := make(chan DebugStep)
	go func() {
		match(runner, []rune(input), debug, 0)
		close(debug)
	}()

	steps := make([]DebugStep, 0)
	for step := range debug {
		steps = append(steps, step)
	}

	return steps
}

func match(r *Runner, input []rune, debugCh chan DebugStep, offset int) bool {
	if debugCh != nil {
		debugCh <- DebugStep{Snapshot: r.snapshot(), CurrentIndex: offset, Offset: offset, Status: r.status()}
	}

	r.Reset()

	for i, char := range input {
		r.Next(char)

		status := r.status()

		if status == StatusFail {
			return match(r, input[1:], debugCh, offset+1)
		}

		step := DebugStep{
			Snapshot:     r.snapshot(),
			Offset:       offset,
			CurrentIndex: offset + i,
			Status:       status,
		}

		step.CurrentIndex++
		// greedy match
		if status == StatusSuccess {
			debugCh <- step
			return true
		}

		if debugCh != nil {
			debugCh <- step
		}
	}

	return r.status() == StatusSuccess
}

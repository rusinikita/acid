package runner

import (
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
)

type StepIterator struct {
	runner   *Runner
	sequence sequence.Sequence
	results  <-chan event.Event
}

func NewIterator(r *Runner, s sequence.Sequence) *StepIterator {
	return &StepIterator{
		runner:   r,
		sequence: s,
		results:  make(chan event.Event),
	}
}

func (i *StepIterator) Run() {
	i.results = i.runner.Run(i.sequence)
}

func (i *StepIterator) Next() (event.Event, bool) {
	e, ok := <-i.results

	return e, ok
}

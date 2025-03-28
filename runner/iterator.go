package runner

import (
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
)

type StepIterator struct {
	runner  *Runner
	results <-chan event.Event
}

func NewIterator(r *Runner) *StepIterator {
	return &StepIterator{
		runner:  r,
		results: make(chan event.Event),
	}
}

func (i *StepIterator) Run(s sequence.Sequence) {
	i.results = i.runner.Run(s)
}

func (i *StepIterator) Next() (event.Event, bool) {
	e, ok := <-i.results

	return e, ok
}

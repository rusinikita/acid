package runner

import (
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
)

type StepIterator struct {
	runner   *Runner
	sequence call.Sequence
	results  <-chan event.Event
}

func (i *StepIterator) Run() {
	i.results = i.runner.Run(i.sequence)
}

func (i *StepIterator) Next() (event.Event, bool) {
	e, ok := <-i.results

	return e, ok
}

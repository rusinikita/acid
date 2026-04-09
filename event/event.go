package event

import (
	"github.com/rusinikita/acid/call"
	"slices"
)

type Event struct {
	trx       call.TrxID
	step      *call.Step
	result    *call.ExecResult
	waiting   []call.TrxID
	testSetup bool
	start     bool
	done      bool
}

func Start() Event { return Event{start: true} }
func Done() Event  { return Event{done: true} }

func (e Event) IsStart() bool { return e.start }
func (e Event) IsDone() bool  { return e.done }

func Call(step call.Step, waiting []call.TrxID) Event {
	return Event{
		trx:       step.Trx,
		step:      &step,
		waiting:   waiting,
		testSetup: step.TestSetup,
	}
}

func Result(step call.Step, result call.ExecResult, waiting []call.TrxID) Event {
	return Event{
		trx:       step.Trx,
		result:    &result,
		waiting:   waiting,
		testSetup: step.TestSetup,
	}
}

func (e Event) Trx() call.TrxID {
	return e.trx
}

func (e Event) Step() *call.Step {
	return e.step
}

func (e Event) Result() *call.ExecResult {
	return e.result
}

func (e Event) IsWaiting(trx call.TrxID) bool {
	return slices.Contains(e.waiting, trx)
}

func (e Event) TestSetup() bool {
	return e.testSetup
}

func (e Event) Waiting() []call.TrxID {
	return e.waiting
}

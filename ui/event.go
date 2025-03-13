package ui

import "github.com/rusinikita/acid/call"

type Event struct {
	trx    call.TrxID
	step   *call.Step
	result *call.ExecResult
}

func Call(step call.Step) Event {
	return Event{
		trx:  step.Trx,
		step: &step,
	}
}

func Result(step call.Step, result call.ExecResult) Event {
	return Event{
		trx:    step.Trx,
		result: &result,
	}
}

func (e Event) IsWaiting() bool {
	return e.step == nil && e.result == nil
}

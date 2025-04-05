package run

import (
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"slices"
)

type eventData struct {
	events        []event.Event
	steps         []call.Step
	transactions  []call.TrxID
	onlyStepsMode bool
}

func (ed *eventData) add(e event.Event) {
	ed.events = append(ed.events, e)

	if step := e.Step(); step != nil {
		ed.steps = append(ed.steps, *step)
	}

	if !slices.Contains(ed.transactions, e.Trx()) {
		ed.transactions = append(ed.transactions, e.Trx())
	}
}

func (ed *eventData) headers() []string {
	headers := make([]string, len(ed.transactions))
	for i, trx := range ed.transactions {
		if trx == "" {
			headers[i] = "No tx requests"
			continue
		}

		headers[i] = "TX - " + string(trx)
	}

	return headers
}

func (ed *eventData) clean() {
	ed.events = ed.events[:0]
	ed.transactions = ed.transactions[:0]
}

func (ed *eventData) At(row, cell int) string {
	trx := ed.transactions[cell]

	if ed.onlyStepsMode {
		return StepStr(ed.steps[row], trx)
	}

	e := ed.events[row]

	s := Cell(e, trx)

	return s
}

func (ed *eventData) Rows() int {
	if ed.onlyStepsMode {
		return len(ed.steps)
	}

	return len(ed.events)
}

func (ed *eventData) Columns() int {
	return len(ed.transactions)
}

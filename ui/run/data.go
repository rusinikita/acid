package run

import (
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"slices"
)

type eventData struct {
	showSetupEvents bool
	onlyStepsMode   bool

	events       []event.Event
	steps        []call.Step
	transactions []call.TrxID

	hasNoTxCalls bool
}

func (ed *eventData) visibleTransactions() []call.TrxID {
	if !ed.showSetupEvents && !ed.hasNoTxCalls && len(ed.transactions) > 1 {
		return ed.transactions[1:]
	}

	return ed.transactions
}

func (ed *eventData) add(e event.Event) {
	if !e.TestSetup() && e.Trx() == "" {
		ed.hasNoTxCalls = true
	}

	ed.events = append(ed.events, e)

	if step := e.Step(); step != nil {
		ed.steps = append(ed.steps, *step)
	}

	if !slices.Contains(ed.transactions, e.Trx()) {
		ed.transactions = append(ed.transactions, e.Trx())
	}
}

func (ed *eventData) headers() []string {
	transactions := ed.visibleTransactions()

	headers := make([]string, len(transactions))
	for i, trx := range transactions {
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
	ed.steps = ed.steps[:0]
	ed.transactions = []call.TrxID{""}
	ed.hasNoTxCalls = false
}

func (ed *eventData) At(row, cell int) string {
	trx := ed.visibleTransactions()[cell]

	if ed.onlyStepsMode {
		step := ed.steps[row]
		if step.TestSetup && !ed.showSetupEvents {
			return ""
		}

		return StepStr(step, trx)
	}

	e := ed.events[row]
	if e.TestSetup() && !ed.showSetupEvents {
		return ""
	}

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
	return len(ed.visibleTransactions())
}

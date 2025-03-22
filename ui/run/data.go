package run

import (
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"slices"
)

type eventData struct {
	events             []event.Event
	transactions       []call.TrxID
	selectedTrxIndex   call.TrxID
	selectedEventIndex int
}

func (ed *eventData) add(e event.Event) {
	ed.events = append(ed.events, e)

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

func (ed *eventData) At(row, cell int) string {
	e := ed.events[row]
	trx := ed.transactions[cell]

	s := e.Cell(trx)

	return s
}

func (ed *eventData) Rows() int {
	return len(ed.events)
}

func (ed *eventData) Columns() int {
	return len(ed.transactions)
}

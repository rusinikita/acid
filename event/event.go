package event

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rusinikita/acid/call"
	"slices"
)

type Event struct {
	trx     call.TrxID
	step    *call.Step
	result  *call.ExecResult
	waiting []call.TrxID
}

func Call(step call.Step, waiting []call.TrxID) Event {
	return Event{
		trx:     step.Trx,
		step:    &step,
		waiting: waiting,
	}
}

func Result(step call.Step, result call.ExecResult, waiting []call.TrxID) Event {
	return Event{
		trx:     step.Trx,
		result:  &result,
		waiting: waiting,
	}
}

func (e Event) Trx() call.TrxID {
	return e.trx
}

func (e Event) Cell(trx call.TrxID) string {
	if e.trx != trx {
		if slices.Contains(e.waiting, trx) {
			return "waiting"
		}

		return ""
	}

	if e.step == nil && e.result == nil {
		return "waiting"
	}

	if e.step != nil {
		if e.step.Code != "" {
			return e.step.Code
		}

		return "Transaction " + e.step.TrxCommand.String()
	}

	result := e.result

	if result.Error != nil {
		return fmt.Sprintln("error:", result.Error)
	}

	if result.Rows == nil {
		return fmt.Sprintln("rows affected:", result.RowsAffected)
	}

	return table.New().Headers(result.Rows.Columns...).Rows(result.Rows.Rows...).String()
}

func (e Event) View() string {
	if e.step != nil {
		return fmt.Sprintln("step", e.step.Trx, e.step.Code)
	}

	result := e.result
	if result == nil {
		return ""
	}

	if result.Error != nil {
		return fmt.Sprintln("error:", result.Error)
	}

	if result.Rows == nil {
		return fmt.Sprintln("rows affected:", result.RowsAffected)
	}

	r := fmt.Sprintln("columns:", result.Rows.Columns)
	for i, row := range result.Rows.Rows {
		r += fmt.Sprintln("row:", i, row)
	}

	return r
}

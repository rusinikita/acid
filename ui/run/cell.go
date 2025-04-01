package run

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/ui/code"
	"github.com/rusinikita/acid/ui/theme"
)

func Cell(e event.Event, trx call.TrxID) string {
	if e.Trx() != trx {
		if e.IsWaiting(trx) {
			return theme.EventTypeStyle.Render("waiting")
		}

		return ""
	}

	step := e.Step()
	result := e.Result()

	if step == nil && e.Result() == nil {
		return theme.EventTypeStyle.Render("waiting")
	}

	if step != nil {
		if step.Code != "" {
			return theme.EventTypeStyle.Render("request") + "\n" + code.Highlight(step.Code)
		}

		return "Transaction " + theme.SQLKeywordStyle.Render(step.TrxCommand.String())
	}

	response := theme.EventTypeStyle.Render("response")

	if result.Error != nil {
		return response + "\n" + theme.ErrorResponseStyle.Render(result.Error.Error())
	}

	if result.Rows == nil {
		return response + "\n" + "rows affected: " + fmt.Sprint(result.RowsAffected)
	}

	return response + "\n" + table.New().Headers(result.Rows.Columns...).Rows(result.Rows.Rows...).String()
}

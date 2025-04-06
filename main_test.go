package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
	"testing"
)

func TestName(t *testing.T) {
	conn := db.Connect()

	r := runner.New(conn)
	i := runner.NewIterator(r)

	var (
		model tea.Model = run.NewRunTable(i, "pg")
		cmd             = model.Init()
	)

	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 20})
	model, _ = model.Update(router.Message{DataInt: 0})

	for {
		if cmd != nil {
			msg := cmd()

			model, cmd = model.Update(msg)
		}

		view := model.View()
		println(len(view))

		if cmd == nil {
			view := model.View()
			println(len(view))

			break
		}
	}
}

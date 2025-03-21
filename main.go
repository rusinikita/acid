package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/ui"
	"log"
)

var (
	tx1 = call.TrxID("first")
	tx2 = call.TrxID("second")

	sequence = call.Sequence{Calls: []call.Step{
		call.Call("CREATE TABLE IF NOT EXISTS exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
		call.Call("TRUNCATE TABLE exec_test"),
		call.Begin(tx1),
		call.Begin(tx2),
		call.Call("insert into exec_test (name) values ('biba')", tx1),
		call.Call("select * from exec_test", tx2),
		call.Commit(tx1),
		call.Call("select * from exec_test", tx2),
		call.Call("select * from exec_test"),
	}}
)

func main() {
	conn := db.Connect()

	r := runner.New(conn)
	i := runner.NewIterator(r, sequence)

	model := ui.NewRunTable(i)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

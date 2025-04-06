package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/list"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
	"log"
	"os"
)

var (
	mainSequence = sequence.Sequences[0]
)

func main() {
	conn := db.Connect()
	driver := os.Getenv("DB_DRIVER")

	r := runner.New(conn)
	i := runner.NewIterator(r)

	routes := map[string]tea.Model{
		"run":  run.NewRunTable(i, driver),
		"list": list.New(driver),
	}

	app := router.NewRouter(routes, router.Route("list"))

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

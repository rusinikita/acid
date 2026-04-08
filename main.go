package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rusinikita/acid/config"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/list"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "path to .toml sequence file")
	flag.StringVar(&filePath, "f", "", "path to .toml sequence file")
	flag.Parse()

	var sequences []sequence.Sequence
	startRoute := router.Route("list")

	if filePath != "" {
		db.LoadEnv(filepath.Dir(filePath))
		seq, err := config.Load(filePath)
		if err != nil {
			log.Fatalf("load %s: %v", filePath, err)
		}
		sequences = []sequence.Sequence{seq}
		startRoute = router.Route("run", 0)
	} else {
		sequences = sequence.Sequences
	}

	conn := db.Connect()
	driver := os.Getenv("DB_DRIVER")

	r := runner.New(conn)
	i := runner.NewIterator(r)

	routes := map[string]tea.Model{
		"run":  run.NewRunTable(i, driver, sequences),
		"list": list.New(driver),
	}

	app := router.NewRouter(routes, startRoute)

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

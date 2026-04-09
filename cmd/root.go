package cmd

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/rusinikita/acid/config"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/list"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
)

var filePath string

var rootCmd = &cobra.Command{
	Use:   "acid",
	Short: "acid – interactive database concurrency tester",
	Run: func(cmd *cobra.Command, args []string) {
		runStandalone()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "path to .toml sequence file")
}

func runStandalone() {
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

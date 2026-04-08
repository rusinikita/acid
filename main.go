package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rusinikita/acid/client"
	"github.com/rusinikita/acid/config"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/server"
	"github.com/rusinikita/acid/ui/list"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			runServe()
			return
		case "run":
			runClient()
			return
		}
	}

	runStandalone()
}

func runStandalone() {
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

func runServe() {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.String("port", "7331", "TCP port to listen on")
	_ = fs.Parse(os.Args[2:])

	addr := ":" + *port
	srv := server.New(addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}

	log.Printf("acid server listening on %s", addr)

	source := runner.NewChannelSource(srv.Channel())

	routes := map[string]tea.Model{
		"run": run.NewServerRunTable(source, *port),
	}

	app := router.NewRouter(routes, router.Route("run"))

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

func runClient() {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	serverAddr := fs.String("server", "localhost:7331", "server address (host:port)")
	filePath := fs.String("file", "", "path to .toml sequence file")
	fs.StringVar(filePath, "f", "", "path to .toml sequence file")
	_ = fs.Parse(os.Args[2:])

	var seq sequence.Sequence

	if *filePath != "" {
		db.LoadEnv(filepath.Dir(*filePath))
		var err error
		seq, err = config.Load(*filePath)
		if err != nil {
			log.Fatalf("load %s: %v", *filePath, err)
		}
	} else {
		args := fs.Args()
		if len(args) == 0 {
			log.Fatal("usage: acid run <sequence-name> [--server host:port]")
		}
		name := args[0]
		var ok bool
		for _, s := range sequence.Sequences {
			if s.Name == name {
				seq, ok = s, true
				break
			}
		}
		if !ok {
			log.Fatalf("sequence %q not found", name)
		}
	}

	conn := db.Connect()
	if err := client.Run(conn, seq, *serverAddr); err != nil {
		log.Fatalf("client: %v", err)
	}
}

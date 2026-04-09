package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/server"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/run"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start acid in server mode, listening for remote sequence runs",
	Run: func(cmd *cobra.Command, args []string) {
		runServe()
	},
}

func init() {
	serveCmd.Flags().StringVar(&port, "port", "7331", "TCP port to listen on")
	rootCmd.AddCommand(serveCmd)
}

func runServe() {
	addr := ":" + port
	srv := server.New(addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}

	log.Printf("acid server listening on %s", addr)

	source := runner.NewChannelSource(srv.Channel())

	routes := map[string]tea.Model{
		"run": run.NewServerRunTable(source, port, srv.ToggleCh()),
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

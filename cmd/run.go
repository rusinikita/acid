package cmd

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rusinikita/acid/client"
	"github.com/rusinikita/acid/config"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/sequence"
)

var serverAddr string
var runFilePath string

var runCmd = &cobra.Command{
	Use:   "run [sequence-name]",
	Short: "Send a sequence to a remote acid server for execution",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runClient(args)
	},
}

func init() {
	runCmd.Flags().StringVar(&serverAddr, "server", "localhost:7331", "server address (host:port)")
	runCmd.Flags().StringVarP(&runFilePath, "file", "f", "", "path to .toml sequence file")
	rootCmd.AddCommand(runCmd)
}

func runClient(args []string) {
	var seq sequence.Sequence

	if runFilePath != "" {
		db.LoadEnv(filepath.Dir(runFilePath))
		var err error
		seq, err = config.Load(runFilePath)
		if err != nil {
			log.Fatalf("load %s: %v", runFilePath, err)
		}
	} else {
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
	if err := client.Run(conn, seq, serverAddr); err != nil {
		log.Fatalf("client: %v", err)
	}
}

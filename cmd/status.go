package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/rusinikita/acid/db"
)

var statusPort string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check database connectivity and acid server reachability",
	Run: func(cmd *cobra.Command, args []string) {
		checkStatus()
	},
}

func init() {
	statusCmd.Flags().StringVarP(&statusPort, "port", "p", "7331", "port of the acid server to check")
	rootCmd.AddCommand(statusCmd)
}

func checkStatus() {
	allOK := true

	if err := db.TryConnect(); err != nil {
		fmt.Printf("database is NOT reachable: %v\n", err)
		allOK = false
	} else {
		fmt.Println("database is reachable")
	}

	url := "http://localhost:" + statusPort + "/health"
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		fmt.Printf("acid server is NOT running on port %s\n", statusPort)
		allOK = false
	} else {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("acid server is running on port %s\n", statusPort)
		} else {
			fmt.Printf("unexpected response from server: %d\n", resp.StatusCode)
			allOK = false
		}
	}

	if !allOK {
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var statusPort string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if an acid server is running on the given port",
	Run: func(cmd *cobra.Command, args []string) {
		checkStatus()
	},
}

func init() {
	statusCmd.Flags().StringVarP(&statusPort, "port", "p", "7331", "port of the acid server to check")
	rootCmd.AddCommand(statusCmd)
}

func checkStatus() {
	url := "http://localhost:" + statusPort + "/health"
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		fmt.Printf("acid server is NOT running on port %s\n", statusPort)
		os.Exit(1)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("acid server is running on port %s\n", statusPort)
	} else {
		fmt.Printf("unexpected response from server: %d\n", resp.StatusCode)
		os.Exit(1)
	}
}

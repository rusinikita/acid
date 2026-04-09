package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var togglePort string

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle result visibility on the running acid server",
	Run: func(cmd *cobra.Command, args []string) {
		toggleMode()
	},
}

func init() {
	toggleCmd.Flags().StringVarP(&togglePort, "port", "p", "7331", "port of the acid server")
	rootCmd.AddCommand(toggleCmd)
}

func toggleMode() {
	url := "http://localhost:" + togglePort + "/toggle-mode"
	resp, err := http.Post(url, "", nil) //nolint:noctx
	if err != nil {
		fmt.Printf("acid server is NOT running on port %s\n", togglePort)
		os.Exit(1)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error: %s\n", body)
		os.Exit(1)
	}

	fmt.Println(string(body))
}

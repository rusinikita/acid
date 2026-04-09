package cmd

import (
	"github.com/spf13/cobra"

	"github.com/rusinikita/acid/initcmd"
)

var initCmd = &cobra.Command{
	Use:   "init [dir]",
	Short: "Scaffold a DB learning environment in the current or given directory",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initcmd.Run(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

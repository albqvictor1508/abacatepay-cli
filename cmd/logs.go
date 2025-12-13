package cmd

import "github.com/spf13/cobra"

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Receive trigger logs",
	RunE:  log,
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func log(cmd *cobra.Command, args []string) error {
}

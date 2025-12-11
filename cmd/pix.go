package cmd

import "github.com/spf13/cobra"

var pixCmd = &cobra.Command{
	Use:   "pix",
	Short: "Manage your PIX operations",
}

func init() {
	rootCmd.AddCommand(pixCmd)
}

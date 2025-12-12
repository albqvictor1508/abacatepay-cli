package cmd

import "github.com/spf13/cobra"

var mrrCmd = &cobra.Command{
	Use:   "mrr",
	Short: "Manage your MRR (Monthly Recurring Revenue)",
}

func init() {
	rootCmd.AddCommand(mrrCmd)
}

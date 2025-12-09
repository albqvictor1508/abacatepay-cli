package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "abacate",
    SilenceUsage:  true,
    SilenceErrors: true,
	Short: "AbacatePay CLI - Interact with the AbacatePay platform",
}

func Exec() {
    cobra.CheckErr(rootCmd.Execute());
}

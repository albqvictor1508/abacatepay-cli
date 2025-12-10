package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "abacate",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "AbacatePay CLI - Interact with the AbacatePay platform",
	Long: `The AbacatePay CLI provides a fast and secure way to interact with the 
AbacatePay platform directly from your terminal.
	
It allows developers to manage API keys, switch profiles, inspect account details, 
trigger test events, handle webhooks, create payments, generate resources, and 
perform common development tasks with a single command.

This CLI was designed to be intuitive, modern and highly scriptable, integrating 
seamlessly with CI/CD pipelines, local development flows, and production environments.`,
}

func Exec() {
	cobra.CheckErr(rootCmd.Execute())
}

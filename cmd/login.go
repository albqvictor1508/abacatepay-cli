package cmd

import (
	"os/exec"

	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/spf13/cobra"
)

var (
	apiKey string
	name   string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in on Abacatepay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return loginOnAbacatepay()
	},
}

func init() {
	loginCmd.Flags().StringVar(&apiKey, "key", "", "API Key do Abacate Pay")
	loginCmd.Flags().StringVar(&name, "name", "", "Nome do profile (mín 3, máx 50 caracteres)")

	rootCmd.AddCommand(loginCmd)
}

func loginOnAbacatepay() error {
	var cmd *exec.Cmd

	if apiKey == "" || name == "" {
		if err := prompts.
	}

	return cmd.Start()
}

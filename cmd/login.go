package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/spf13/cobra"
	"honnef.co/go/tools/config"
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
	loginCmd.Flags().StringVar(&apiKey, "key", "", "Abacate Pay's API Key")
	loginCmd.Flags().StringVar(&name, "name", "", "Profile name (min 3, max 50 characters)")

	rootCmd.AddCommand(loginCmd)
}

func loginOnAbacatepay() error {
	var cmd *exec.Cmd

	if apiKey == "" || name == "" {
		if err := prompts.InteractiveMode(apiKey); err != nil {
			return err
		}
	}

	if len(name) < 3 || len(name) > 50 {
		return fmt.Errorf("the profile name must to be at 3 and 50 characters")
	}

	cfg, err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error to load config: %w", err)
	}

	return cmd.Start()
}

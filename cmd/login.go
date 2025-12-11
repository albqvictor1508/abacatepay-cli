package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/albqvictor1508/abacatepay-cli/internal/config"
	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in on Abacatepay with an API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logIn()
	},
}

var (
	key  string
	name string
)

func init() {
	loginCmd.Flags().StringVar(&key, "key", "", "Abacate Pay's API Key")
	loginCmd.Flags().StringVar(&name, "name", "", "Name for the profile (Min 3, Max 50 chars.)")

	rootCmd.AddCommand(loginCmd)
}

func logIn() error {
	var cmd *exec.Cmd

	if key == "" {
		input, err := prompts.AskAPIKey()
		if err != nil {
			return err
		}

		key = input
	}

	if name == "" {
		input, err := prompts.AskProfileName()
		if err != nil {
			return err
		}

		name = input
	}

	config, err := config.Load()

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unknown error while loading the profile (%w)", err)
	}

	if config.Exists(name) {
		return fmt.Errorf("Name \"%s\" is already being used as a profile", name)
	}

	if err := config.Add(name, key); err != nil {
		return err
	}

	return cmd.Start()
}

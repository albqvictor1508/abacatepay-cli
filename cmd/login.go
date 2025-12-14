package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/albqvictor1508/abacatepay-cli/internal/api"
	"github.com/albqvictor1508/abacatepay-cli/internal/config"
	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
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

	cfg, err := config.Load()

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unknown error while loading the profile (%w)", err)
	}

	if !utils.IsOnline() {
		return fmt.Errorf("Unable to verify API key. You're offline.")
	}

	// TODO: adicionar um spinner ou algo do tipo aq pra ficar daorinha

	valid, err := api.ValidateAPIKey(key)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("Invalid API key.")
	}

	if cfg.Exists(name) {
		return fmt.Errorf("Name \"%s\" is already being used as a profile", name)
	}

	if err := cfg.Add(name, key); err != nil {
		return err
	}

	fmt.Printf("Successfully logged in. Profile '%s' is now active.\n", name)

	return cmd.Start()
}

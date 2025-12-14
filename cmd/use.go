package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albqvictor1508/abacatepay-cli/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Switch to a different profile",
	Args:  cobra.ExactArgs(1),
	RunE:  use,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func use(_ *cobra.Command, args []string) error {
	var cmd *exec.Cmd

	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if !cfg.ProfileExists(profileName) {
		return fmt.Errorf("profile '%s' does not exist", profileName)
	}

	if err := cfg.SetCurrent(profileName); err != nil {
		return fmt.Errorf("error setting current profile: %w", err)
	}

	fmt.Printf("Switched to profile: '%s'\n", profileName)
	return cmd.Start()
}

package cmd

import (
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Abacatepy CLI when new version is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() error {
	latest, err := CheckForUpdate()
	if err != nil {
		log.Println("error checking for update: %w", err)
	}
}

func CheckForUpdate() (*selfupdate.Release, error) {
	currentVersion, err := semver.Parse(version)

	latest, found, err := selfupdate.DetectLatest("albqvictor1508/abacatepay-cli")
	if err != nil {
		return nil, err
	}

	if !found || latest.Version.LTE(currentVersion) {
		return nil, nil
	}

	return latest, nil
}

func ShowUpdate(version string) {
}

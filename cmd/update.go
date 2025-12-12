package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Abacatepy CLI when new version is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		update()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() {
	latest, err := CheckForUpdate()
	if err != nil {
		log.Println("error checking for update: %w", err)
		return
	}

	if latest == nil {
		fmt.Println("Current version is the latest") // TODO: trocar pelo surveys
		return
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		if os.IsPermission(err) {
			log.Println("permission denied, please run the update command with sudo: sudo abacate update")
			return
		}

		log.Println("error ocurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
	return
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
	latest, err := CheckForUpdate()

	if err != nil || latest == nil {
		return
	}

	// TODO: colocar um popup aq, tipo "Update Available" pra quando o usuario
	// iniciar a CLI aparecer
}

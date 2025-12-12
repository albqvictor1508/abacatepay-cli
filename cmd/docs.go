package cmd

import (
	"github.com/spf13/cobra"
	"os/exec"
	"runtime"
)

var useReference bool

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Open AbacatePay documentation in the browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://docs.abacatepay.com"

		if useReference {
			url += "/api-reference"
		}

		return open(url)
	},
}

func init() {
	docsCmd.Flags().BoolVarP(&useReference, "reference", "r", false, "Go directly to the API reference")

	rootCmd.AddCommand(docsCmd)
}

func open(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

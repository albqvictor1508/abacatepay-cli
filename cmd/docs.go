package cmd

import (
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Open AbacatePay documentation in the browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		return openBrowser("https://docs.abacatepay.com")
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}

func openBrowser(url string) error {
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

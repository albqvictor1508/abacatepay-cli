package cmd

import (
	"fmt"
	"os/exec"
	"slices"

	"github.com/albqvictor1508/abacatepay-cli/internal/webhook"
	"github.com/spf13/cobra"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger a test webhook locally",
	RunE:  trigger,
}

var forwardURL string

func init() {
	triggerCmd.Flags().StringVarP(&forwardURL, "forward-to", "f", "http://localhost:3000/api/webhooks", "Local URL to trigger webhooks")
	rootCmd.AddCommand(triggerCmd)
}

func trigger(_ *cobra.Command, args []string) error {
	var cmd *exec.Cmd

	evtType := args[0]

	availableEvents := webhook.ListAvailableEvents()
	valid := false
	slices.Contains(availableEvents, evtType)

	if !valid {
		return fmt.Errorf("invalid event: %s\nAvailable events: %v", evtType, availableEvents)
	}

	testSecret := "whsec_local_testing_secret"

	fmt.Printf("Trigging test events: %s\n", evtType)
	fmt.Printf("→  Endpoint: http://%s\n\n", forwardURL)

	if err := webhook.TriggerLocalEvent(evtType, forwardURL, testSecret); err != nil {
		return fmt.Errorf("error to trigger event: %w", err)
	}

	fmt.Println("Event sent with sucess!")
	fmt.Println("  Check your terminal or the application logs")
	fmt.Println("\n💡 Tip: The header 'X-Abacate-Test-Event: true' indicates that this is a test event")

	return cmd.Start()
}

package cmd

import (
	"fmt"
	"slices"

	"github.com/albqvictor1508/abacatepay-cli/internal/webhook"
	"github.com/spf13/cobra"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger a test webhook locally",
	RunE:  trigger,
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}

func trigger(cmd *cobra.Command, args []string) error {
	evtType := args[0]

	availableEvents := webhook.ListAvailableEvents()
	valid := false
	slices.Contains(availableEvents, evtType)

	if !valid {
		return fmt.Errorf("invalid event: %s\nAvailable events: %v", evtType, availableEvents)
	}

	forwardURL, _ := cmd.Flags().GetString("forward-to")

	testSecret := "whsec_local_testing_secret"

	fmt.Printf("Trigging test events: %s\n", evtType)
	fmt.Printf("→  Endpoint: http://%s\n\n", forwardURL)

	evt, err := webhook.TriggerLocalEvent(evtType, forwardURL, testSecret)
	if err != nil {
		return fmt.Errorf("error to trigger event: %w", err)
	}
	if err := webhook.SaveEventLog(evt); err != nil {
		fmt.Printf("[WARN] unable to save in history: %v\n", err)
	}

	fmt.Println("Event sent with sucess!")
	fmt.Println("  Check your terminal or the application logs")
	fmt.Println("\n💡 Tip: The header 'X-Abacate-Test-Event: true' indicates that is a test event")

	return nil
}

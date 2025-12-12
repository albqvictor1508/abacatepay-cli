package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
	"github.com/albqvictor1508/abacatepay-cli/internal/webhook"
	"github.com/spf13/cobra"
)

var forwardTo string

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen webhooks and make forward to localhost",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listen()
	},
}

func init() {
	listenCmd.Flags().StringVarP(&forwardTo, "forward-to", "f", "localhost:3000/webhook", "Local URL to receive webhooks")

	rootCmd.AddCommand(listenCmd)
}

func listen() error {
	apiKey, err := utils.GetAPIKey()
	if err != nil {
		return fmt.Errorf("error to get API key: %w\nUse 'abacate login' first", err)
	}

	proxy := webhook.NewProxy(apiKey, forwardTo)

	fmt.Println("Abacate Pay")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("→  Forward to: http://%s\n", forwardTo)
	fmt.Printf("→  Webhook Signing Secret: %s\n", proxy.SigningSecret())

	if err := proxy.Connect(); err != nil {
		return fmt.Errorf("error to connect: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		errChan <- proxy.Listen()
	}()

	select {
	case <-sigChan:
		fmt.Println("\n\n Finalizing webhook proxy...")
		proxy.Close()
		return nil
	case err := <-errChan:
		return fmt.Errorf("error in proxy: %w", err)
	}
}

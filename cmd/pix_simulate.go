package cmd

import (
	"fmt"

	"github.com/albqvictor1508/abacatepay-cli/internal/api"
	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use: "simulate",
	Args: cobra.ExactArgs(1),
	Short: "Simulates payment for a QRCode Pix",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := utils.GetAPIKey()

		if err != nil {
			return err
		}

		data, err := api.NewClient(key).Request(api.RequestOptions{
			Method: "POST",
			Route: "/pixQrCode/simulate-payment?id="+args[0],
		})

		if err != nil {
			return err
		}

		brCode, _ := data["brCode"].(string)

		fmt.Printf("BrCode \"%s\"", brCode)

		return nil
	},
}

func init() {
	pixCmd.AddCommand(simulateCmd)
}

package cmd

import (
	"fmt"

	"github.com/albqvictor1508/abacatepay-cli/internal/api"
	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
	"github.com/spf13/cobra"
)

var pixStatusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.ExactArgs(1),
	Short: "Check QRCode Pix payment status",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := utils.GetAPIKey()

		if err != nil {
			return err
		}

		id := args[0]

		client := api.NewClient(key)

		data, err := client.Request(api.RequestOptions{
			Key:    key,
			Method: "GET",
			Route:  "/pixQrCode/check?id=" + id,
		})

		if err != nil {
			return err
		}

		status, _ := data["status"].(string)

		fmt.Printf("Your PIX \"%s\" status is \"%s\"\n", id, status)

		return nil
	},
}

func init() {
	pixCmd.AddCommand(pixStatusCmd)
}

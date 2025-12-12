package cmd

import (
	"github.com/albqvictor1508/abacatepay-cli/internal/api"
	"github.com/albqvictor1508/abacatepay-cli/internal/ui"
	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Retrieve your account details including balance information",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := utils.GetAPIKey()

		if err != nil {
			return err
		}

		data, err := api.NewClient(key).Request(api.RequestOptions{
			Method: "GET",
			Route:  "/store/get",
		})

		if err != nil {
			return err
		}

		return ui.PrintJSON(data)
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
}

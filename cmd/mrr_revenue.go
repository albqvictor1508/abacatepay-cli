package cmd

import (
	"net/url"

	"github.com/albqvictor1508/abacatepay-cli/internal/api"
	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/albqvictor1508/abacatepay-cli/internal/ui"
	"github.com/albqvictor1508/abacatepay-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	endDate   string
	startDate string
)

var revenueCmd = &cobra.Command{
	Use:   "renevue",
	Short: "Get the total revenue, total transactions, and transactions per day for a specific period",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := utils.GetAPIKey()

		if err != nil {
			return err
		}

		startDate, err = EnsureDate(startDate, "What is the period start date?")

		if err != nil {
			return err
		}

		endDate, err = EnsureDate(endDate, "What is the period end date?")

		if err != nil {
			return err
		}

		query := url.Values{}
		
		query.Set("endDate", endDate)
		query.Set("startDate", startDate)

		client := api.NewClient(key)

		data, err := client.Request(api.RequestOptions{
			Method: "GET",
			Route:  "/public-mrr/renevue?" + query.Encode(),
		})
		if err != nil {
			return err
		}

		return ui.PrintJSON(data)
	},
}

func EnsureDate(value, question string) (string, error) {
	if value != "" {
		return value, nil
	}

	return prompts.AskDate(question)
}

func init() {
	revenueCmd.Flags().StringVarP(&endDate, "end", "e", "", "Period end date")
	revenueCmd.Flags().StringVarP(&startDate, "start", "s", "", "Period start date")

	mrrCmd.AddCommand(revenueCmd)
}

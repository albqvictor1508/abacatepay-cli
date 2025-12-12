package prompts

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

func AskDate(message string) (string, error) {
	var date string

	for {
		err := survey.AskOne(&survey.Input{
			Message: message,
		}, &date, survey.WithValidator(survey.Required))

		if err != nil {
			return "", err
		}

		_, perr := time.Parse("2006-01-02", date)

		if perr == nil {
			return date, nil
		}

		fmt.Println("Unknown date format, use YYYY-MM-DD")
	}
}

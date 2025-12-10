package prompts

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func AskAPIKey() (string, error) {
	var useFromEnv bool
	var key = os.Getenv("ABACATE_PAY_API_KEY")

	isUsingFromEnv := len(key) > 0

	if isUsingFromEnv {
		err := survey.AskOne(&survey.Confirm{
			Default: false,
			Message: "We detected an API key in ABACATE_PAY_API_KEY enviroment. Do you want to use it?",
		}, &useFromEnv)

		if err != nil {
			return "", err
		}
	}

	if !useFromEnv {
		err := survey.AskOne(&survey.Password{
			Message: "What's the API key?",
		}, &key, survey.WithValidator(survey.Required))

		if err != nil {
			return "", err
		}
	}

	return key, nil
}

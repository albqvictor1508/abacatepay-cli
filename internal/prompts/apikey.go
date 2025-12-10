package prompts

import "github.com/AlecAivazis/survey/v2"

func AskAPIKey() (string, error) {
	var key string

	err := survey.AskOne(&survey.Password{
		Message: "What's the API key?",
	}, &key, survey.WithValidator(survey.Required))

	return key, err
}

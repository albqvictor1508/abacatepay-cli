package prompts

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func InteractiveMode(apiKey string) error {
	envKey := os.Getenv("ABACATE_PAY_API_KEY")

	if envKey != "" {
		useEnv := false
		prompt := &survey.Confirm{
			Message: "Foi encontrada uma API key no ambiente (ABACATE_PAY_API_KEY). Deseja usá-la?",
			Default: true,
		}
		if err := survey.AskOne(prompt, &useEnv); err != nil {
			return err
		}
		if useEnv {
			apiKey = envKey
		}
	}

	if apiKey == "" {
		prompt := &survey.Password{
			Message: "Qual a sua API key do Abacate Pay?",
		}
		if err := survey.AskOne(prompt, &apiKey, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	return nil
}

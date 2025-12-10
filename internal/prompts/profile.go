package prompts

import "github.com/AlecAivazis/survey/v2"

func AskProfileName() (string, error) {
	var name string

	err := survey.AskOne(&survey.Input{
		Message: "How we should name it?",
	}, &name, survey.WithValidator(survey.MinLength(3)))

	return name, err
}

package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func AskProfileName() (string, error) {
	var name string

	const maxNameLen = 50

	err := survey.AskOne(&survey.Input{
		Message: "How we should name it?",
	}, &name, survey.WithValidator(survey.MinLength(3)))

	if len(name) > maxNameLen {
		return "", fmt.Errorf("Expecting a name with length between 3 and %v (Got %v)", maxNameLen, len(name))
	}

	return name, err
}

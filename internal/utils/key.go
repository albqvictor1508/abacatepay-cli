package utils

import (
	"errors"

	"github.com/albqvictor1508/abacatepay-cli/internal/config"
	"github.com/albqvictor1508/abacatepay-cli/internal/prompts"
	"github.com/zalando/go-keyring"
)

func GetAPIKey() (string, error) {
	key, err := config.GetKeyring("current")

	if err == nil {
		return key, nil
	}

	if !errors.Is(err, keyring.ErrNotFound) {
		return "", err
	}

	return prompts.AskAPIKey()
}

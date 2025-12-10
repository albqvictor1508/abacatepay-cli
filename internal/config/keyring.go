package config

import (
	keyring "github.com/zalando/go-keyring"
)

var service = "abacate"

func SaveToKeyring(name, apiKey string) error {
	return keyring.Set(service, name, apiKey)
}

func GetFromKeyring(name string) (string, error) {
	return keyring.Get(service, name)
}

func DeleteFromKeyring(name string) error {
	return keyring.Delete(service, name)
}

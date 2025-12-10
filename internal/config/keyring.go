package config

import (
	keyring "github.com/zalando/go-keyring"
)

const service = "abacate"

func SaveKeyring(name, key string) error {
	return keyring.Set(service, name, key)
}

func GetKeyring(name string) (string, error) {
	return keyring.Get(service, name)
}

func DeleteKeyring(name string) error {
	return keyring.Delete(service, name)
}

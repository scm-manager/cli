package auth

import (
	"fmt"
	"github.com/zalando/go-keyring"
)

func ReadApiKey(username string) (string, error) {
	apiKey, err := keyring.Get("scm-cli", username)
	if err != nil {
		return apiKey, fmt.Errorf("could not create cli config file: %w", err)
	}
	return apiKey, nil
}

func StoreApiKey(username string, apiKey string) error {
	err := keyring.Set("scm-cli", username, apiKey)
	if err != nil {
		return fmt.Errorf("could not create cli config file: %w", err)
	}
	return nil
}

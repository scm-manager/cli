package auth

import (
	"github.com/zalando/go-keyring"
	"log"
)

func ReadApiKey(username string) string {
	apiKey, err := keyring.Get("scm-cli", username)
	if err != nil {
		log.Fatal("Could not create cli config file", err)
	}
	return apiKey
}

func StoreApiKey(username string, apiKey string) {
	err := keyring.Set("scm-cli", username, apiKey)
	if err != nil {
		log.Fatal("Could not create cli config file", err)
	}
}

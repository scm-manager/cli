package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type CreateApiKeyRequest struct {
	Hostname string `json:"hostname"`
}

var KeyName = "scm-cli"

func Create(serverUrl string, username string, password string) (string, error) {
	// Create api key on server
	hostname, err := resolveHostname()
	if err != nil {
		return "", err
	}
	loginRequest := CreateApiKeyRequest{Hostname: hostname}
	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode(loginRequest)
	if err != nil {
		return "", fmt.Errorf("could not encode hostname: %w", err)
	}
	req, _ := http.NewRequest("POST", serverUrl+"/api/v2/cli/login", payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send login request: %w", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read login response: %w", err)
	}
	return string(body), nil
}

func Remove(apiKeyName string, serverUrl string, username string, apiKey string) error {
	// Remove api key on server
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode("")
	hostname, err := resolveHostname()
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("DELETE", serverUrl+"/api/v2/cli/logout/"+hostname, payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not revoke api key on server: %w", err)
	}
	defer res.Body.Close()
	return nil
}

func resolveHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not resolve hostname:  %w", err)

	}
	return hostname, nil
}

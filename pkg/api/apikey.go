package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/scm-manager/cli/pkg"
	"io"
	"io/ioutil"
	"net/http"
)

type CreateApiKeyRequest struct {
	ApiKey string `json:"apiKey"`
}

func Create(serverUrl string, username string, password string, apiKeyName string) (string, error) {
	loginRequest := CreateApiKeyRequest{ApiKey: apiKeyName}
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(loginRequest)
	if err != nil {
		return "", fmt.Errorf("could not encode hostname: %w", err)
	}
	req, _ := http.NewRequest("POST", serverUrl+"/api/v2/cli/login", payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	client := pkg.CreateHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send login request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode == 401 {
		return "", createLoginError("Login failed. Please check username and password.", res)
	} else if res.StatusCode == 404 {
		return "", createLoginError("Failed to create API key. Please make sure that API keys are enabled in the global configuration.", res)
	} else if res.StatusCode == 403 {
		return "", createLoginError("You do not have the permission to use API keys. This is mandatory for the usage of the cli client.", res)
	} else if res.StatusCode == 409 {
		return "", createLoginError("There already seems to be an api key for the CLI client issued for this user/system. Please delete this api key manually on your scm server and retry to log in via the CLI.", res)
	} else if res.StatusCode >= 400 {
		return "", createLoginError("Could not create new api key on server.", res)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read login response: %w", err)
	}
	return string(body), nil
}

func Remove(serverUrl string, apiKey string, apiKeyName string) error {
	// Remove api key on server
	req, err := http.NewRequest("DELETE", serverUrl+"/api/v2/cli/logout/"+apiKeyName, nil)
	if err != nil {
		return fmt.Errorf("could not create delete request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)
	client := pkg.CreateHttpClient()
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not revoke api key on server: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode >= 400 {
		return fmt.Errorf("could not remove api key. Server returned status code: %d", res.StatusCode)
	}
	return nil
}

func createLoginError(message string, res *http.Response) error {
	return fmt.Errorf("Server returned status code: %d\n%v", res.StatusCode, message)
}

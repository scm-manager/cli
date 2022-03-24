package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zalando/go-keyring"
	"net/http"
	"os"
)

func Logout() error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	// Send logout request to server
	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode("")
	hostname, err := ResolveHostname()
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("DELETE", config.ServerUrl+"/api/v2/cli/logout/"+hostname, payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	apiKey, err := ReadApiKey(config.Username)
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.Username, apiKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not revoke api key on server: %w", err)
	}
	defer res.Body.Close()

	// Delete local config
	filePath, err := ResolveConfigFilePath()
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("could not delete config file: %w", err)
	}

	// Delete stored api key
	err = keyring.Delete("scm-cli", config.Username)
	if err != nil {
		if err != nil {
			return fmt.Errorf("could not delete stored api key: %w", err)
		}
	}
	return nil
}

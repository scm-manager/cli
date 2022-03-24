package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Login() error {
	// Collect login parameters
	serverUrl := os.Args[2]
	username, password, err := ReadCredentials()
	if err != nil {
		return err
	}

	// Send login request
	hostname, err := ResolveHostname()
	if err != nil {
		return err
	}
	loginRequest := LoginRequest{Hostname: hostname}
	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode(loginRequest)
	if err != nil {
		return fmt.Errorf("could not encode hostname: %w", err)
	}
	req, _ := http.NewRequest("POST", serverUrl+"/api/v2/cli/login", payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send login request: %w", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read login response: %w", err)
	}

	// Create .scm-cli.json config file
	configFilePath, err := ResolveConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("could not create cli config file: %w", err)
	}
	configuration := Configuration{ServerUrl: serverUrl, ApiKey: string(body), Username: username}
	jsonConfig, err := json.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("could not marshal config json: %w", err)
	}
	_, err = file.Write(jsonConfig)
	if err != nil {
		return fmt.Errorf("could not write config to file: %w", err)

	}

	// Store api key to keyring
	err = StoreApiKey(username, configuration.ApiKey)
	if err != nil {
		return err
	}
	return nil
}

func ReadCredentials() (string, string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := r.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("could not read username: %w", err)
	}
	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", fmt.Errorf("could not read password: %w", err)
	}
	return strings.TrimSpace(username), string(password), nil
}

type LoginRequest struct {
	Hostname string `json:"hostname"`
}

package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ReadConfig() (*Configuration, error) {
	configFilePath, err := ResolveConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("could not resolve config file path: %w", err)
	}

	_, err = os.Stat(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("could not find config file: %w", err)
		}
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	config := &Configuration{}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)

	}
	key, err := ReadApiKey(config.Username)
	if err != nil {
		return config, err
	}
	config.ApiKey = key
	return config, err
}

func ResolveConfigFilePath() (string, error) {
	// Read config (server url, username) => .scm-cli.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return path.Join(homeDir, ".scm-cli.json"), nil
}

func ResolveHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not resolve hostname:  %w", err)

	}
	return hostname, nil
}

type Configuration struct {
	ServerUrl string
	Username  string
	ApiKey    string `json:"-"`
}
